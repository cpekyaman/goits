# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

# app specific parameters
BINARY_PATH=../bin
BINARY_NAME=gdevl
ENTRY_POINT=./main/main.go
BINARY_CMD=$(BINARY_PATH)/$(BINARY_NAME).exe
PKG_ROOT=github.com/cpekyaman/devtools

# simple general tasks
all: clean test build
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_CMD)
test: 
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_CMD) $(ENTRY_POINT)