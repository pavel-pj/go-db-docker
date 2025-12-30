MIGRATIONS_DIR = database/migrations

.PHONY: migrate-up migrate-down migrate-create

migrate-up:
	docker run --rm \
		-v $(PWD)/$(MIGRATIONS_DIR):/migrations \
		-v $(PWD):/data \
		migrate/migrate:v4.18.2 \
		-path=/migrations \
		-database="sqlite:///data/$(DB_PATH)?x-no-tx-wrap=true" \
		up

migrate-down:
	docker run --rm \
		-v $(PWD)/$(MIGRATIONS_DIR):/migrations \
		-v $(PWD):/data \
		migrate/migrate:v4.18.2 \
		-path=/migrations \
		-database="sqlite:///data/$(DB_PATH)?x-no-tx-wrap=true" \
		down 1

migrate-create:
	@read -p "Enter migration name: " name; \
	docker run --rm \
		-v $(PWD)/$(MIGRATIONS_DIR):/migrations \
		migrate/migrate:v4.18.2 \
		create -ext sql -dir /migrations -seq $$name

main:
	docker compose exec backend go run main.go		