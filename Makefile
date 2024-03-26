NAME = serenity
COMMIT = $(shell git rev-parse --short HEAD)
TAG = $(shell git describe --tags --always)
VERSION = $(TAG:v%=%)

PARAMS = -v -trimpath -ldflags "-X 'github.com/sagernet/serenity/constant.Version=$(VERSION)' -s -w -buildid="
MAIN_PARAMS = $(PARAMS)
MAIN = ./cmd/serenity
PREFIX ?= $(shell go env GOPATH)

.PHONY: release docs

build:
	go build $(MAIN_PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(NAME) $(MAIN_PARAMS) $(MAIN)

fmt:
	@gofumpt -l -w .
	@gofmt -s -w .
	@gci write --custom-order -s standard -s "prefix(github.com/sagernet/)" -s "default" .

fmt_install:
	go install -v mvdan.cc/gofumpt@latest
	go install -v github.com/daixiang0/gci@latest

lint:
	GOOS=linux golangci-lint run ./...
	GOOS=android golangci-lint run ./...
	GOOS=windows golangci-lint run ./...
	GOOS=darwin golangci-lint run ./...
	GOOS=freebsd golangci-lint run ./...

lint_install:
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest

release:
	goreleaser release --clean --skip publish
	mkdir dist/release
	mv dist/*.tar.gz dist/*.deb dist/*.rpm dist/*.pkg.tar.zst dist/release
	ghr --replace --draft --prerelease -p 3 "v${VERSION}" dist/release
	rm -r dist/release

release_repo:
	goreleaser release -f .goreleaser.fury.yaml --clean

release_install:
	go install -v github.com/goreleaser/goreleaser@latest
	go install -v github.com/tcnksm/ghr@latest

docs:
	mkdocs serve

publish_docs:
	mkdocs gh-deploy -m "Update" --force --ignore-version --no-history

docs_install:
	pip install --force-reinstall mkdocs-material=="9.*" mkdocs-static-i18n=="1.2.*"

clean:
	rm -rf bin dist serenity
	rm -f $(shell go env GOPATH)/serenity

update:
	git fetch
	git reset FETCH_HEAD --hard
	git clean -fdx