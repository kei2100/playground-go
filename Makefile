PACKAGES := $(shell go list ./...)

# setup tasks
.PHONY: setup
setup:
	go install golang.org/x/tools/cmd/goimports@latest

# development tasks
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: test
test:
	go test -race $(PACKAGES)

.PHONY: test.nocache
test.nocache:
	go test -count=1 -race $(PACKAGES)
