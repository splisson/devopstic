# go option
GO        ?= go
TAGS      :=
TESTS     := .
TESTFLAGS :=
LDFLAGS   := -w -s
GOFLAGS   :=
BINDIR    := $(CURDIR)/bin
VERSION   := 0.2

# Required for globs to work correctly
SHELL=/bin/bash

all: devopstic

dependencies:
	dep ensure

devopstic: dependencies
	GOBIN=$(BINDIR) $(GO) install $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' main/main.go
	mv bin/main bin/devopstic

docker:
	docker build -t splisson/devopstic:v$(VERSION) .
	docker push splisson/devopstic:v$(VERSION)
