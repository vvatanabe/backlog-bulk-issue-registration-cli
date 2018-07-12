NAME = backlog-bulk-issue-registration
PKG = github.com/vvatanabe/backlog-bulk-issue-registration-cli
COMMIT = $$(git describe --tags --always)
DATE = $$(date '+%Y-%m-%d_%H:%M:%S')
BUILD_LDFLAGS = -X $(PKG)/internal.commit=$(COMMIT) -X $(PKG)/internal.date=$(DATE)
RELEASE_BUILD_LDFLAGS = -s -w $(BUILD_LDFLAGS)

.PHONY: build
build:
	go build -ldflags="$(BUILD_LDFLAGS)" -o ./$(NAME) ./cmd/main.go

.PHONY: test
test:
	go test -v -race ./internal/...

.PHONY: devel-deps
devel-deps:
	go get github.com/mattn/goveralls
	go get github.com/motemen/gobump/cmd/gobump
	go get github.com/Songmu/ghch/cmd/ghch
	go get github.com/Songmu/goxz/cmd/goxz
	go get github.com/tcnksm/ghr

.PHONY: crossbuild
crossbuild: devel-deps
	$(eval ver = $(shell gobump show -r))
	goxz -pv=$(ver) -arch=386,amd64 -build-ldflags="$(RELEASE_BUILD_LDFLAGS)" \
	  -o=$(NAME) -d=./dist/$(ver) ./cmd