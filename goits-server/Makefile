# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOGEN=$(GOCMD) generate

# app specific parameters
BINARY_PATH=../bin
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
frameworkTest: ormTest validationTest cachingTest routingTest

# test tasks for application part
projectTest:
	$(GOTEST) -v $(PKG_ROOT)/application/project
applicationTest: projectTest

# mock generation for framework components
svcMock:
	$(GOGEN) "github.com/cpekyaman/goits/framework/services"
validationMock:
	$(GOGEN) "github.com/cpekyaman/goits/framework/valiation"
cacheMock:
	$(GOGEN) "github.com/cpekyaman/goits/framework/caching"

# running application or application tasks
run-db-migrate: build
	$(BINARY_CMD) migrate up
run-db-rollback: build
	$(BINARY_CMD) migrate down
run-server: build
	$(BINARY_CMD) server start