GO ?= go
GOBIN ?= $$($(GO) env GOPATH)/bin

.PHONY: all
all: build vet test binary

.PHONY: binary
binary: CGO_ENABLED=0
binary: dist FORCE
	$(GO) version
	$(GO) build -race -trimpath -ldflags "$(LDFLAGS)" -o dist/loggy .

dist:
	mkdir $@

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test -v -failfast ./...

.PHONY: build
build: CGO_ENABLED=0
build:
	$(GO) build -race -trimpath -ldflags "$(LDFLAGS)" ./...

.PHONY: clean
clean:
	$(GO) clean
	rm -rf dist/

FORCE:
