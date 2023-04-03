NAME = serenity
COMMIT = $(shell git rev-parse --short HEAD)
TAGS ?= with_acme

PARAMS = -v -trimpath -tags "$(TAGS)" -ldflags "-s -w -buildid="
MAIN = ./cmd/serenity
PREFIX ?= $(shell go env GOPATH)

.PHONY: test release

build:
	go build $(PARAMS) $(MAIN)

install:
	go build -o $(PREFIX)/bin/$(NAME) $(PARAMS) $(MAIN)

fmt:
	@gofumpt -l -w .
	@gofmt -s -w .
	@gci write --custom-order -s "standard,prefix(github.com/sagernet/),default" .

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
	goreleaser release --rm-dist --skip-publish || exit 1
	mkdir dist/release
	mv dist/*.tar.gz dist/*.zip dist/*.deb dist/*.rpm dist/release
	ghr --delete --draft --prerelease -p 3 $(shell git describe --tags) dist/release
	rm -r dist

release_install:
	go install -v github.com/goreleaser/goreleaser@latest
	go install -v github.com/tcnksm/ghr@latest
	clean:
	rm -rf bin dist serenity
	rm -f $(shell go env GOPATH)/serenity

update:
	git fetch
	git reset FETCH_HEAD --hard
	git clean -fdx
