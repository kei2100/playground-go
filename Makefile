.PHONY: test fmt vet lint setup vendor

PACKAGES := $(shell go list ./...)
DIRS := $(shell go list -f '{{.Dir}}' ./...)

setup:
	@which dep > /dev/null 2>&1 || go get -u github.com/golang/dep/cmd/dep
	@which golint > /dev/null 2>&1 || go get -u github.com/golang/lint/golint
	@which goimports > /dev/null 2>&1 || go get -u golang.org/x/tools/cmd/goimports
	@which richgo > /dev/null 2>&1 || go get -u github.com/kyoh86/richgo

vendor: setup vendor/.timestamp

vendor/.timestamp: $(shell find $(DIRS) -name '*.go')
	dep ensure -v
	touch vendor/.timestamp

vet:
	go vet $(PACKAGES)

lint: setup
	! find $(DIRS) -name '*.go' | xargs goimports -d | grep '^'
	echo $(PACKAGES) | xargs -n 1 golint -set_exit_status

fmt: setup
	find $(DIRS) -name '*.go' | xargs goimports -w

test: setup vendor
	go test -v -race $(PACKAGES) | richgo testfilter
