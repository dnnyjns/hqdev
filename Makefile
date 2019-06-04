# Go parameters
GOCMD=go
GOBUILD=${GOCMD} build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=hqdev

all: build

build:
				$(GOBUILD) -o $(BINARY_NAME) -v

clean:
				$(GOCLEAN)
				rm -f $(BINARY_NAME)

run: build
				./$(BINARY_NAME) run

