# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# app specific parameters
BINARY_PATH=..
BINARY_NAME=goits
ENTRY_POINT=./main/main.go
BINARY_CMD=$(BINARY_PATH)/$(BINARY_NAME).exe
PKG_ROOT=github.com/cpekyaman/goits

# simple general tasks
all: clean test build
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_CMD)
test: 
	$(GOTEST) -v ./...
build:
	$(GOBUILD) -o $(BINARY_CMD) $(ENTRY_POINT)

# test tasks for framework part
ormTest:
	$(GOTEST) -v $(PKG_ROOT)/framework/orm
validationTest:
	$(GOTEST) -v $(PKG_ROOT)/framework/validation
cachingTest:
	$(GOTEST) -v $(PKG_ROOT)/framework/caching
routingTest:
	$(GOTEST) -v $(PKG_ROOT)/framework/routing
frameworkTest: ormTest validationTest cachingTest

# test tasks for application part
projectTest:
	$(GOTEST) -v $(PKG_ROOT)/application/project
applicationTest: projectTest


# running application or application tasks
run-db-migrate:
	$(BINARY_CMD) migrate up
run-db-rollback:
	$(BINARY_CMD) migrate down
run-server:
	$(BINARY_CMD) server start