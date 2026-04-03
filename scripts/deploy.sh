#!/bin/bash
# =============================================================
# deploy.sh — Deploy / redeploy ทั้ง stack
# รัน: bash scripts/deploy.sh
# =============================================================
set -e

echo "==> [1/4] Pull latest code (ถ้าใช้ git)"
# git pull origin main   # uncomment ถ้า push code ขึ้น EC2 ผ่าน git

echo "==> [2/4] Check .env.prod exists"
if [ ! -f .env.prod ]; then
  echo "ERROR: .env.prod not found!"
  echo "Copy .env.prod.example -> .env.prod แล้วแก้ค่าให้ครบก่อน"
  exit 1
fi

echo "==> [3/4] Build & start containers"
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build

echo "==> [4/4] Run database migrations"
sleep 5   # รอ postgres พร้อม
docker compose -f docker-compose.prod.yml exec backend ./main migrate 2>/dev/null || true

echo ""
echo "==> Deploy complete!"
docker compose -f docker-compose.prod.yml ps
