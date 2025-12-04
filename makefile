build:
	@go build -o tmp/main main.go

run: build
	@./tmp/main

run-mock:
	@go run ./mock-server/mock-server.go
