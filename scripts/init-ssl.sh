#!/bin/bash
# =============================================================
# init-ssl.sh — ขอ SSL certificate ครั้งแรกจาก Let's Encrypt
# รัน: bash scripts/init-ssl.sh
# =============================================================
set -e

DOMAIN="checkjop-mathcom-cu.com"
EMAIL="your-email@example.com"   # <-- เปลี่ยนเป็น email จริง

# Step 1: ใช้ config แบบ HTTP-only ก่อน (ยังไม่มี cert)
echo "==> Switching to HTTP-only nginx config for cert challenge..."
cd nginx/conf.d
mv checkjop.conf checkjop.conf.bak
mv checkjop-init.conf checkjop.conf
cd ../..

# Step 2: Start nginx ด้วย config ชั่วคราว
echo "==> Starting nginx..."
docker compose -f docker-compose.prod.yml up -d nginx frontend backend postgres

# Step 3: ขอ certificate
echo "==> Requesting SSL certificate..."
docker compose -f docker-compose.prod.yml run --rm certbot certonly \
  --webroot \
  --webroot-path=/var/www/certbot \
  --email "$EMAIL" \
  --agree-tos \
  --no-eff-email \
  -d "$DOMAIN" \
  -d "www.$DOMAIN"

# Step 4: สลับกลับเป็น HTTPS config
echo "==> Switching back to HTTPS nginx config..."
cd nginx/conf.d
mv checkjop.conf checkjop-init.conf
mv checkjop.conf.bak checkjop.conf
cd ../..

# Step 5: Reload nginx
echo "==> Reloading nginx..."
docker compose -f docker-compose.prod.yml exec nginx nginx -s reload

echo ""
echo "==> SSL setup complete! Site is live at https://$DOMAIN"
