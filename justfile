# run dev
dev *flags:
  go run main.go {{ flags }}

# build main app
build:
	@echo "--- building the program ---"
	GO111MODULE=on go mod tidy
	CGO_ENABLED=0 go build -o bin/main main.go
	@echo "--- build complete ---"

# run compiled app
run *flags:
	./bin/main {{ flags }}

# run test
test:
	@echo "--- running tests ---"
	go test ./... -v
	@echo "--- tests complete ---"
