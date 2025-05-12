build:
	@echo "--- building the program ---"
	GO111MODULE=on go mod tidy
	CGO_ENABLED=0 go build -o bin/main main.go
	@echo "--- build complete ---"

run:
	@bin/main
