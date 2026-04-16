.PHONY: build test lint clean run

BIN := driftwatch
CMD := ./cmd/driftwatch

build:
	go build -o $(BIN) $(CMD)

test:
	go test ./... -v -count=1

lint:
	golangci-lint run ./...

clean:
	rm -f $(BIN)

run: build
	./$(BIN) $(ARGS)

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report written to coverage.html"
