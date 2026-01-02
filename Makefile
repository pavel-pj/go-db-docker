.PHONY: migrate migrate-create migrate-up migrate-down migrate-status db-shell
 
#Загружаем переменные из .env
include .env
export

MIGRATIONS_DIR=./backend/database/migrations
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
	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" up
# Откатить миграции
migrate-down:
	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" down

# Показать статус
migrate-status:
	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" version

# Откатить 1 миграцию
migrate-rollback:
	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" down 1
main:
	docker compose exec backend go run main.go		

