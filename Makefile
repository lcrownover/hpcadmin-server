build:
	@go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

run: build
	@./bin/hpcadmin-server
