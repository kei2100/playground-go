.PHONY: test fmt vet lint

PACKAGES := $(shell go list ./...)
DIRS := $(shell go list -f '{{.Dir}}' ./...)

test:
	go test -v -race $(PACKAGES)

fmt:
	find $(DIRS) -name '*.go' | xargs goimports -w

vet:
	go vet $(PACKAGES)

lint:
	echo $(PACKAGES) | xargs -n 1 golint -set_exit_status
