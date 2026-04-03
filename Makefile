# =============================================================
# Makefile — shortcuts สำหรับ deploy checkjop
# =============================================================

.PHONY: help up down restart logs deploy ssl seed-data

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

seed-data:
	@echo "==> Importing curriculum..."
	@curl -s -X POST http://localhost:8080/api/v1/import/curriculum-csv \
		-F "file=@checkjop-be/csv_data/version9/example_checkjop - curriculum.csv" | cat
	@echo "\n==> Importing categories..."
	@curl -s -X POST http://localhost:8080/api/v1/import/category-csv \
		-F "file=@checkjop-be/csv_data/version9/example_checkjop - catagory-2.csv" | cat
	@echo "\n==> Importing courses 2566..."
	@curl -s -X POST http://localhost:8080/api/v1/import/course-csv-with-year \
		-F "file=@checkjop-be/csv_data/version9/example_checkjop - course_Present_2566-6.csv" \
		-F "year=2566" | cat
	@echo "\n==> Importing courses 2567..."
	@curl -s -X POST http://localhost:8080/api/v1/import/course-csv-with-year \
		-F "file=@checkjop-be/csv_data/version9/example_checkjop - course_Present_2567-5.csv" \
		-F "year=2567" | cat
	@echo "\n==> Importing courses 2568..."
	@curl -s -X POST http://localhost:8080/api/v1/import/course-csv-with-year \
		-F "file=@checkjop-be/csv_data/version9/example_checkjop - course_Present_2568-5.csv" \
		-F "year=2568" | cat
	@echo "\n==> Done!"

clean:
	docker compose -f docker-compose.prod.yml down -v --remove-orphans
