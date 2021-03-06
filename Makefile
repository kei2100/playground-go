PACKAGES := $(shell go list ./...)

GREP := grep
ifeq ($(OS),Windows_NT)
	GREP := findstr
endif

# setup tasks
.PHONY: setup
setup:
	go get golang.org/x/tools/cmd/goimports
	go get golang.org/x/lint/golint
	go get github.com/kyoh86/richgo
	go mod tidy

# development tasks
.PHONY: fmt
fmt:
	goimports -w .

.PHONY: lint
lint:
	goimports -d . | $(GREP) "^" && exit 1 || exit 0
	golint -set_exit_status $(PACKAGES)

.PHONY: vet
vet:
	go vet $(PACKAGES)

.PHONY: test
test:
	richgo test -v -race $(PACKAGES)

.PHONY: test.nocache
test.nocache:
	richgo test -count=1 -v -race $(PACKAGES)
