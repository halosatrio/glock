build:
	@GO111MODULE=on go mod tidy
	@CGO_ENABLED=0 go build -o bin/main

run:
	@bin/main
