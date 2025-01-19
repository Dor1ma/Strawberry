build:
	go build -o strawberry.exe ./cmd/strawberry/main.go

test:
	go test ./...
