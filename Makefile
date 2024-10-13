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
	@go mod tidy && go build -ldflags="-X 'github.com/mstyushin/$(APP_NAME)/pkg/config.Version=$(VERSION)' -X 'github.com/mstyushin/$(APP_NAME)/pkg/config.Hash=$(HASHCOMMIT)' -X 'github.com/mstyushin/$(APP_NAME)/pkg/config.VersionDate=$(VERSIONDATE)'" -o bin/$(APP_NAME) github.com/mstyushin/$(APP_NAME)/cmd/server
	@chmod +x bin/$(APP_NAME)

run: build
	@mkdir -p bin log
	@bin/$(APP_NAME) > log/$(APP_NAME).log 2>&1 & echo "$$!" > /tmp/$(APP_NAME).pid

e2e-test:
	@go mod tidy && go test -v ./...

stop:
	-pkill $(APP_NAME)

docker-image:
	docker build . -t  mstyushin/$(APP_NAME):$(VERSION) -t mstyushin/$(APP_NAME):latest

docker-push:
	docker push mstyushin/$(APP_NAME):$(VERSION)
	docker push mstyushin/$(APP_NAME):latest
