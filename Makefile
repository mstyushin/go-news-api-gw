SHELL = /bin/bash

APP_NAME := go-news-api-gw

$(eval TAGVERSION := $(shell git describe --tags))
$(eval HASHCOMMIT := $(shell git log --pretty=tformat:"%h" -n1 ))
$(eval BRANCHNAME := $(shell git branch --show-current))
ifeq ($(TAGVERSION),undefined)
    # default tag is undefined
    VERSION := $(BRANCHNAME)
else ifeq ($(TAGVERSION),)
    # is empty tag 
    VERSION := $(BRANCHNAME)
else
    VERSION := $(TAGVERSION)
endif
$(eval VERSIONDATE := $(shell git show -s --format=%cI $($VERSION)))

clean: stop
	@rm -f bin/*
	@rm -f log/*
	@rm -f /tmp/$(APP_NAME).pid

build:
	@go mod tidy && go build -ldflags="-X 'github.com/mstyushin/go-news-api-gw/pkg/config.Version=$(VERSION)' -X 'github.com/mstyushin/go-news-api-gw/pkg/config.Hash=$(HASHCOMMIT)' -X 'github.com/mstyushin/go-news-api-gw/pkg/config.VersionDate=$(VERSIONDATE)'" -o bin/$(APP_NAME) github.com/mstyushin/go-news-api-gw/cmd/server
	@chmod +x bin/$(APP_NAME)

run: build
	@mkdir -p bin log
	@bin/$(APP_NAME) > log/$(APP_NAME).log 2>&1 & echo "$$!" > /tmp/$(APP_NAME).pid

test:
	@go mod tidy && go test -v ./...

stop:
	-pkill $(APP_NAME)
