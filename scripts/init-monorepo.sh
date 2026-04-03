#!/bin/bash
# =============================================================
# init-monorepo.sh — รันบนเครื่อง local ครั้งเดียว
# เพื่อ init monorepo จาก fe+be ที่มีอยู่
# =============================================================
set -e

cd /Users/get/Documents/coding/checkjop

echo "==> Remove nested .git folders (fe + be มี git ของตัวเอง)"
rm -rf checkjop-fe/.git
rm -rf checkjop-be/.git

echo "==> Remove old CI from be (ย้ายมาไว้ที่ monorepo แล้ว)"
rm -rf checkjop-be/.github

echo "==> Init git monorepo"
git init
git add .
git commit -m "chore: init monorepo (fe + be + deploy configs)"

echo ""
echo "==> ต่อไปทำบน GitHub:"
echo "1. สร้าง repo ใหม่ที่ github.com/cu-devclub/checkjop (อย่า init README)"
echo "2. รันคำสั่งนี้:"
echo ""
echo "   git remote add origin https://github.com/cu-devclub/checkjop.git"
echo "   git branch -M main"
echo "   git push -u origin main"
