.PHONY: migrate migrate-create migrate-up migrate-down migrate-status db-shell
 
#Загружаем переменные из .env
include .env
export

MIGRATIONS_DIR=./backend/database/migrations
DB_URL_DOCKER=postgres://$(DB_USER):$(DB_PASSWORD)@postgres:5432/$(DB_NAME)?sslmode=disable
 
#goose

goose-create:
	@echo "migration name :"
	@read -p "> " name && \
	docker compose exec backend goose -dir /app/database/migrations create $$name sql

goose-status:
	docker compose exec backend goose -dir database/migrations status
goose-up:
	docker compose exec backend goose -dir database/migrations up	
goose-rollback:
	docker compose exec backend goose -dir database/migrations down	

main:
	docker compose exec backend go run main.go		
bash:
	docker compose exec backend sh

right:
	 sudo chown -R $$USER:$$USER ./

 
#для migration
#migrate-up:
#	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" up
# Откатить миграции
#migrate-down:
#	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" down

# Показать статус
#migrate-status:
#	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" version

# Откатить 1 миграцию
#migrate-rollback:
#	docker compose exec backend migrate -path /app/database/migrations -database "$(DB_URL_DOCKER)" down 1
 
