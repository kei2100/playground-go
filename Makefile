#GO ?= go1.21rc2
GO ?= go

PACKAGES := $(shell $(GO) list ./...)

# setup tasks
.PHONY: setup
setup:
#	## use 1.21rc2
#	go install golang.org/dl/go1.21rc2@latest && go1.21rc2 download
	$(GO) install golang.org/x/tools/cmd/goimports@latest

# development tasks
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: vet
vet:
	$(GO) vet $(PACKAGES)

.PHONY: test
test:
	$(GO) test -race $(PACKAGES)

.PHONY: test.nocache
test.nocache:
	$(GO) test -count=1 -race $(PACKAGES)
