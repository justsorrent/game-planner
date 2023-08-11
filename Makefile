build:
	@go build -o target/game-planner
run-planner:
	@make build && ./target/game-planner
test:
	@go test -v