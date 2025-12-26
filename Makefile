.PHONY: migrate migrate-create migrate-up migrate-down migrate-status db-shell

#DB_URL=postgres://${DB_USER}:${DB_PASSWORD}@localhost:5450/${DB_NAME}?sslmode=disable
#DB_URL_DOCKER=postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable

# Загружаем переменные из .env
include .env
export

MIGRATIONS_DIR=./backend/migrations
DB_URL_DOCKER=postgres://$(DB_USER):$(DB_PASSWORD)@postgres:5432/$(DB_NAME)?sslmode=disable

# Создать файлы миграции НА ХОСТЕ
migrate-create:
	@read -p "Enter migration name: " name; \
	last=$$(ls $(MIGRATIONS_DIR)/*.up.sql 2>/dev/null | wc -l); \
	next=$$(printf "%03d" $$((last + 1))); \
	echo "Creating migration $${next}_$${name}..."; \
	touch $(MIGRATIONS_DIR)/$${next}_$${name}.up.sql; \
	touch $(MIGRATIONS_DIR)/$${next}_$${name}.down.sql; \
	echo "Created: $(MIGRATIONS_DIR)/$${next}_$${name}.{up,down}.sql"
# Запустить миграции
migrate-up:
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" up
# Откатить миграции
migrate-down:
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" down

# Показать статус
migrate-status:
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" version

# Откатить 1 миграцию
migrate-rollback:
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" down 1

# Перезапустить (down + up)
migrate-refresh:
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" down && \
	docker compose exec backend migrate -path /app/migrations -database "$(DB_URL_DOCKER)" up

# Shell БД
db-shell:
	docker compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}		