.PHONY: all fmt test build clean
GOBIN   ?= `which go`
PACKAGE ?= .
GOPATH  ?= $(PWD)

all: fmt test build

build:
	$(GOBIN) build $(PACKAGE)

test:
	$(GOBIN) test $(PACKAGE)/...

clean:
	$(GOBIN) clean $(PACKAGE)
	rm -rv ./pkg/*

fmt:
	$(GOBIN) fmt $(PACKAGE)/...
