build:
	go build -o lox.exe ./cmd/lox/main.go

test:
	go test ./...
