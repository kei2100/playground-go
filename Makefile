GO ?= go

PACKAGES := $(shell $(GO) list ./...)

# setup tasks
.PHONY: setup

setup:
	$(GO) install golang.org/x/tools/cmd/goimports@latest
	$(GO) install github.com/rakyll/gotest@latest

# development tasks
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: test
test:
	gotest -race $(PACKAGES)

.PHONY: test.nocache
test.nocache:
	gotest -count=1 -race $(PACKAGES)
