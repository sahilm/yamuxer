.PHONY: all
all: tools lint test

PKGS := $(shell go list ./... | grep -v /vendor)
.PHONY: test
test: tools
	go test $(PKGS)

sources = $(shell find . -name '*.go')
.PHONY: goimports
goimports: tools
	goimports -w $(sources)

.PHONY: lint
lint: tools
	gometalinter ./... --enable=goimports --enable=gosimple \
	--enable=unparam --enable=unused --disable=gotype --disable=golint -t

BIN_DIR := $(GOPATH)/bin
GOIMPORTS := $(BIN_DIR)/goimports
GOMETALINTER := $(BIN_DIR)/gometalinter
DEP := $(BIN_DIR)/dep

$(GOIMPORTS):
	go get -u golang.org/x/tools/cmd/goimports

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install &> /dev/null

tools: $(GOIMPORTS) $(GOMETALINTER)
