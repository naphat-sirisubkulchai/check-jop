#!/bin/bash
# =============================================================
# setup-ec2.sh — รันครั้งแรกบน EC2 เพื่อติดตั้ง dependencies
# รัน: bash setup-ec2.sh
# =============================================================
set -e

echo "==> [1/5] Update packages"
sudo apt-get update -y && sudo apt-get upgrade -y

echo "==> [2/5] Install Docker"
sudo apt-get install -y ca-certificates curl gnupg lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] \
  https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update -y
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

echo "==> [3/5] Add ubuntu user to docker group"
sudo usermod -aG docker ubuntu
newgrp docker

echo "==> [4/5] Enable Docker on startup"
sudo systemctl enable docker
sudo systemctl start docker

echo "==> [5/5] Add swap (1GB) — สำคัญมากสำหรับ t2.micro"
if [ ! -f /swapfile ]; then
  sudo fallocate -l 1G /swapfile
  sudo chmod 600 /swapfile
  sudo mkswap /swapfile
  sudo swapon /swapfile
  echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
  echo "Swap created"
else
  echo "Swap already exists"
fi

echo ""
echo "==> Done! Log out and back in for docker group to take effect"
echo "==> Then run: bash deploy.sh"
