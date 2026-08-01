#!/usr/bin/env bash
set -euo pipefail

if [[ -n "${ADMIN_EMAIL:-}" && -n "${ADMIN_PASSWORD:-}" ]]; then
  psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    --set=admin_email="$ADMIN_EMAIL" \
    --set=admin_password="$ADMIN_PASSWORD" <<'SQL'
INSERT INTO users (email, password_hash)
VALUES (:'admin_email', crypt(:'admin_password', gen_salt('bf')))
ON CONFLICT (email) DO UPDATE
SET password_hash = EXCLUDED.password_hash;
SQL
else
  echo "ADMIN_EMAIL or ADMIN_PASSWORD not set; skipping initial admin user"
fi
