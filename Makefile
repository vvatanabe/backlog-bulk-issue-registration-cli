.PHONY: deps

NAME = backlog-bulk-issue-registration-cli

default: test

test:
	docker run --rm \
      -v ${PWD}:/go/src/github.com/vvatanabe/backlog-bulk-issue-registration-cli \
      backlog-bulk-issue-registration-cli bash -c ' \
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./... \
      '

build-test-image:
	docker build \
      -t backlog-bulk-issue-registration-cli .