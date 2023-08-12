include .env

build:
	@go build -o target/game-planner
run-planner:
	@make build && ./target/game-planner
test:
	@go test -v
goose-up:
	@cd sql/schema && goose postgres $(DB_STRING) up && cd ../..
goose-down:
	@cd sql/schema && goose postgres $(DB_STRING) down && cd ../..
goose-status:
	@cd sql/schema && goose postgres $(DB_STRING) status && cd ../..