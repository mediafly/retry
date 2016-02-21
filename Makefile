.PHONY: default deps help test

default: help

help:
	@echo "tagets:"
	@echo
	@echo "  deps - install dependencies"
	@echo "  test - run tests"
	@echo

deps:
	@go get github.com/mediafly/math
	@go get github.com/stretchr/testify/assert

test:
	@go test -v .
