.PHONY: test.nocache test fmt vet lint setup vendor

PACKAGES := $(shell go list ./...)
DIRS := $(shell go list -f '{{.Dir}}' ./...)

setup:
	which dep > /dev/null 2>&1 || go get -u github.com/golang/dep/cmd/dep
	which goimports > /dev/null 2>&1 || go get -u golang.org/x/tools/cmd/goimports
	which golint > /dev/null 2>&1 || go get -u golang.org/x/lint/golint
	which richgo > /dev/null 2>&1 || go get -u github.com/kyoh86/richgo

vendor: vendor/.timestamp

vendor/.timestamp: $(shell find $(DIRS) -maxdepth 1 -name '*.go')
	dep ensure -v
	touch vendor/.timestamp

vet:
	go vet $(PACKAGES)

lint:
	! find $(DIRS) -maxdepth 1 -name '*.go' | xargs goimports -d | grep '^'
	echo $(PACKAGES) | xargs -n 1 golint -set_exit_status

fmt:
	find $(DIRS) -maxdepth 1 -name '*.go' | xargs goimports -w

test:
	richgo test -v -race $(PACKAGES)

test.nocache:
	richgo test -count=1 -v -race $(PACKAGES)

