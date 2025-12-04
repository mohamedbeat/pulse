build:
	@go build -o tmp/main main.go

run: build
	@./tmp/main
