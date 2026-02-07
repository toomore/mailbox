#!/bin/bash
# Local checks: go vet, test coverage (requires MariaDB - run sh ./dev-run-mariadb.sh first)
# ref. https://gist.github.com/hailiang/0f22736320abe6be71ce

set -e

go vet ./...
go test -coverprofile=coverage.out -covermode=atomic ./...
go tool cover -func=coverage.out
