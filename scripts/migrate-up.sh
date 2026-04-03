#!/bin/bash
set -e

DB_USER="${DB_USER:-checkjop-admin-user-11501112}"
DB_NAME="${DB_NAME:-checkjop}"

echo "==> Running SQL migrations..."
for f in $(ls checkjop-be/migrations/*.up.sql | sort); do
  echo "  Applying $f..."
  cat "$f" | docker compose -f docker-compose.prod.yml --env-file .env.prod exec -T postgres psql -U "$DB_USER" -d "$DB_NAME"
done
echo "==> Migrations complete."
