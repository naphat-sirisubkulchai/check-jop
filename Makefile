# =============================================================
# Makefile — shortcuts สำหรับ deploy checkjop
# =============================================================

.PHONY: help up down restart logs deploy ssl

help:
	@echo "Available commands:"
	@echo "  make up       — Start all services"
	@echo "  make down     — Stop all services"
	@echo "  make restart  — Restart all services"
	@echo "  make logs     — Tail logs (all services)"
	@echo "  make deploy   — Build & deploy (full)"
	@echo "  make ssl      — Init SSL certificate (ครั้งแรก)"
	@echo "  make ps       — Show running containers"

up:
	docker compose -f docker-compose.prod.yml --env-file .env.prod up -d

down:
	docker compose -f docker-compose.prod.yml down

restart:
	docker compose -f docker-compose.prod.yml --env-file .env.prod restart

logs:
	docker compose -f docker-compose.prod.yml logs -f --tail=100

logs-be:
	docker compose -f docker-compose.prod.yml logs -f --tail=100 backend

logs-fe:
	docker compose -f docker-compose.prod.yml logs -f --tail=100 frontend

ps:
	docker compose -f docker-compose.prod.yml ps

deploy:
	bash scripts/deploy.sh

ssl:
	bash scripts/init-ssl.sh

clean:
	docker compose -f docker-compose.prod.yml down -v --remove-orphans
