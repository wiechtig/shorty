start-dev: stop-dev
	docker compose up -d

stop-dev:
	docker compose down

start-app: stop-app
	docker compose -f compose.app.yml up -d --build

stop-app:
	docker compose -f compose.app.yml down

run:
	SHORTY_DATABASE_URL="postgres://postgres:postgres@localhost:5432/shorty?sslmode=disable" go run ./cmd/shorty

.PHONY: k6
k6:
	k6 run k6/resolver.ts

.PHONY: test
test:
	@echo "Running tests for $(APP_NAME)..."
	go test -v ./...

generate-db:
	sqlc generate
