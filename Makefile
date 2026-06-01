.PHONY: all linux mac mac-arm windows clean

APP_NAME=MONITORING-TOOL
MAIN_PATH=./cmd/agent

all: linux mac mac-arm windows

linux:
	@echo "Building Linux Binary..."
	GOOS=linux GOARCH=amd64 go build -o dist/linux/$(APP_NAME) $(MAIN_PATH)

mac:
	@echo "Building macOS Intel Binary..."
	GOOS=darwin GOARCH=amd64 go build -o dist/mac/$(APP_NAME) $(MAIN_PATH)

mac-arm:
	@echo "Building macOS ARM Binary..."
	GOOS=darwin GOARCH=arm64 go build -o dist/mac-arm/$(APP_NAME) $(MAIN_PATH)

windows:
	@echo "Building Windows Binary..."
	GOOS=windows GOARCH=amd64 go build -o dist/windows/$(APP_NAME).exe $(MAIN_PATH)

clean:
	@echo "Cleaning dist directory..."
	rm -rf dist/
