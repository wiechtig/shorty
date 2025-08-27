start-dev: stop-dev
	docker compose up -d

stop-dev:
	docker compose down

start-app: stop-app
	docker compose -f compose.app.yml up -d --build

stop-app:
	docker compose -f compose.app.yml down

run:
	SHORTY_DATABASE_URL="postgres://postgres:postgres@localhost:5432/shorty?sslmode=disable" \
	SHORTY_OIDC_ISSUER="http://localhost:8080/realms/shorty" \
	SHORTY_OIDC_CLIENT_ID="shorty" \
	go run ./cmd/shorty

.PHONY: hurl
hurl:
	hurl --variable host=http://localhost:8080 --variable shorty=http://localhost:4444 -v ./hurl/basic.hurl

get-token:
	curl -X POST -s "http://localhost:8080/realms/shorty/protocol/openid-connect/token" \
		-H "Content-Type: application/x-www-form-urlencoded" \
		-d "client_id=shorty" \
		-d "username=user@test.com" \
		-d "password=Test123!" \
		-d "grant_type=password" \
		-d "scope=openid profile email" \
		| jq -r '.access_token'

.PHONY: k6
k6:
	k6 run k6/resolver.ts

.PHONY: test
test:
	@echo "Running tests for $(APP_NAME)..."
	go test -v ./...

generate-api:
	go generate ./...

generate-db:
	sqlc generate

generate: generate-api generate-db
