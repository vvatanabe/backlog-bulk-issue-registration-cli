.PHONY: deps

NAME = backlog-bulk-issue-register

default: test

test:
	docker run --rm \
      -v ${PWD}:/go/src/github.com/vvatanabe/backlog-bulk-issue-register \
      backlog-bulk-issue-register bash -c ' \
        go test -v -race -coverprofile=coverage.out -covermode=atomic ./... \
      '

build-test-image:
	docker build \
      -t backlog-bulk-issue-register .