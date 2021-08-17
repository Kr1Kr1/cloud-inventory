OUT_FOLDER := out
BUILDDIR := build
PROJECTNAME := cloud-inventory
PKG := $(PROJECTNAME)
GOPATH := $(shell go env GOPATH)
VERSION := $(shell git describe --abbrev=0 --tags)
GIT_COMMIT_HASH := $(shell git rev-parse HEAD)
CTIMEVAR=-X ${PKG}/utils.GITCOMMIT=${GIT_COMMIT_HASH} -X ${PKG}/utils.VERSION=${VERSION} -X ${PKG}/utils.EXECUTABLE=${PROJECTNAME}
GO_LDFLAGS=-ldflags "-w ${CTIMEVAR}"
GO_LDFLAGS_STATIC=-ldflags "-w ${CTIMEVAR} -extldflags -static"
GOLANGCI_LINT_CMD := golangci-lint

all: static
	@./${PROJECTNAME} -h

run: ## go run main.go
	@go run main.go

run-static: static ## compile & run static executable
	@./${PROJECTNAME}

build: ## Build a dynamic executable or package (not recommended)
	go build -i -v -o ${PROJECTNAME} ${GO_LDFLAGS} ${PKG}

static: ## Build a static executable
	@echo "Compiling static ${PROJECTNAME}..."
	@CGO_ENABLED=0 go build -mod=vendor -o ${PROJECTNAME} -tags "static_build netgo" -installsuffix netgo ${GO_LDFLAGS_STATIC} ${PKG}

vendor: ## Updates the vendoring directory.
	go mod init $(PKG) || true
	go mod tidy
	go mod vendor

clean: ## Remove previous build
	-@rm -f ${PROJECTNAME}

clean-branch: 
	git fetch -p && git branch --merged master | grep -v ^\\* | xargs git branch -d

style: ## Run golangci-lint
	${GOLANGCI_LINT_CMD} run

fmt: ## Verifies all files have been `gofmt`ed.
	@if [[ ! -z "$(shell gofmt -s -l . | grep -v '.pb.go:' | grep -v '.twirp.go:' | grep -v vendor | tee /dev/stderr)" ]]; then \
		exit 1; \
	fi

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: go build clean static style vendor run
