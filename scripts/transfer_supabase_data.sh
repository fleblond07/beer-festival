#!/usr/bin/env bash
set -euo pipefail

if [[ -z "${SUPABASE_DATABASE_URL:-}" ]]; then
  echo "SUPABASE_DATABASE_URL is required" >&2
  exit 1
fi

TARGET_DATABASE_URL="${DATABASE_URL:-postgres://beer_festival:beer_festival@localhost:5432/beer_festival?sslmode=disable}"
DUMP_FILE="${DUMP_FILE:-/tmp/beer-festival-supabase-data.sql}"
for binary in pg_dump psql; do
  if ! command -v "$binary" >/dev/null 2>&1; then
    echo "$binary is required" >&2
    exit 1
  fi
done

pg_dump "$SUPABASE_DATABASE_URL" \
  --schema=public \
  --data-only \
  --column-inserts \
  --table=public.festivals \
  --table=public.breweries \
  --table=public.festivals_breweries \
  --file="$DUMP_FILE"

psql "$TARGET_DATABASE_URL" <<SQL
TRUNCATE TABLE festivals_breweries, festivals, breweries RESTART IDENTITY CASCADE;
SQL

psql "$TARGET_DATABASE_URL" --file="$DUMP_FILE"

psql "$TARGET_DATABASE_URL" <<'SQL'
SELECT setval(pg_get_serial_sequence('festivals', 'id'), COALESCE((SELECT MAX(id) FROM festivals), 1), (SELECT COUNT(*) > 0 FROM festivals));
SELECT setval(pg_get_serial_sequence('breweries', 'id'), COALESCE((SELECT MAX(id) FROM breweries), 1), (SELECT COUNT(*) > 0 FROM breweries));
SQL

if [[ -n "${ADMIN_EMAIL:-}" && -n "${ADMIN_PASSWORD:-}" ]]; then
  psql "$TARGET_DATABASE_URL" \
    --set=admin_email="$ADMIN_EMAIL" \
    --set=admin_password="$ADMIN_PASSWORD" <<'SQL'
INSERT INTO users (email, password_hash)
VALUES (:'admin_email', crypt(:'admin_password', gen_salt('bf')))
ON CONFLICT (email) DO UPDATE
SET password_hash = EXCLUDED.password_hash;
SQL
fi

echo "Transferred Supabase public table data into $TARGET_DATABASE_URL"
