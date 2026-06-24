BINARY=fugu
BUILD_DIR=bin

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/fugu

build-all:
	mkdir -p $(BUILD_DIR)/linux $(BUILD_DIR)/macos $(BUILD_DIR)/windows

	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(BINARY) ./cmd/fugu
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/macos/$(BINARY) ./cmd/fugu
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(BINARY).exe ./cmd/fugu

run:
	go run ./cmd/fugu

test:
	go test ./... -v

fmt:
	go fmt ./...

clean:
	rm -rf $(BUILD_DIR)

gen:
	go run ./cmd/generator/action/map/main.go
	go run ./cmd/generator/action/slice/main.go