lint:
	golangci-lint run ./...
gosec:
	gosec ./...
	govulncheck ./...
run:
	go run main.go
all: lint gosec run

