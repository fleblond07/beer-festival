#!/usr/bin/env bash
set -euo pipefail

if [[ -n "${ADMIN_EMAIL:-}" && -n "${ADMIN_PASSWORD:-}" ]]; then
  export ADMIN_EMAIL ADMIN_PASSWORD
  psql --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    --set ON_ERROR_STOP=1 \
    --file=/docker-entrypoint-initdb.d/admin-user-upsert.psql
else
  echo "ADMIN_EMAIL or ADMIN_PASSWORD not set; skipping initial admin user"
fi
