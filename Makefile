build:
	@go build -o target/game-planner
run-planner:
	@make build && ./target/game-planner
test:
	@go test -v
goose-up:
	@cd sql/schema && goose postgres "user=lucasorrentino dbname=game-planner sslmode=disable" up && cd ../..
goose-down:
	@cd sql/schema && goose postgres "user=lucasorrentino dbname=game-planner sslmode=disable" down && cd ../..
goose-status:
	@cd sql/schema && goose postgres "user=lucasorrentino dbname=game-planner sslmode=disable" status && cd ../..