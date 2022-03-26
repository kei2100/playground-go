GO ?= go

PACKAGES := $(shell $(GO) list ./...)

# setup tasks
.PHONY: setup
setup:
	## use go1.18beta1
	#go get golang.org/dl/go1.18beta1 && go1.18beta1 download
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
