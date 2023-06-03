ifndef $(GOPATH)
	GOPATH=$(shell go env GOPATH)
	export GOPATH
endif

build:
	@go build -o bin/hpcadmin-server cmd/hpcadmin-server/main.go

run: build
	@./bin/hpcadmin-server

docs: build
	@./bin/hpcadmin-server -docs=markdown
