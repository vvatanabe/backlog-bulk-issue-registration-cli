NAME = bbir
PKG = github.com/vvatanabe/backlog-bulk-issue-registration-cli
VERSION = $(shell gobump show -r)
COMMIT = $$(git describe --tags --always)
DATE = $$(date '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X $(PKG)/internal.commit=$(COMMIT) -X $(PKG)/internal.date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

ifeq ($(update),yes)
  u=-u
endif

.PHONY: deps
deps:
	go get ${u} github.com/golang/dep/cmd/dep
	dep ensure

.PHONY: devel-deps
devel-deps: deps
	go get ${u} github.com/mattn/goveralls
	go get ${u} github.com/golang/lint/golint
	go get ${u} github.com/motemen/gobump/cmd/gobump
	go get ${u} github.com/Songmu/ghch/cmd/ghch
	go get ${u} github.com/Songmu/goxz/cmd/goxz
	go get ${u} github.com/tcnksm/ghr

.PHONY: test
test: deps
	go test -v -race -covermode=atomic -coverprofile=coverage.out ./bbir/...

.PHONY: cover
cover: devel-deps
	goveralls -coverprofile=coverage.out -service=travis-ci

.PHONY: lint
lint: devel-deps
	go vet ./bbir/...
	golint -set_exit_status ./bbir/...

.PHONY: bump
bump: devel-deps
	./_tools/bump

.PHONY: build
build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o ./dist/$(NAME) ./cmd/main.go

.PHONY: crossbuild
crossbuild: devel-deps
	$(eval ver = $(shell gobump show -r))
	goxz -pv=$(ver) -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -o=$(NAME) -d=./dist/$(ver) ./cmd

.PHONY: upload
upload:
	ghr -username vvatanabe -replace ${VERSION} ./dist/${ver}

.PHONY: release
release:
	bump crossbuild upload