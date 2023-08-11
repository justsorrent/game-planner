build:
	@go build -o target/game-planner -v
run-planner:
	@make build && ./target/game-planner
test:
	@go test -v