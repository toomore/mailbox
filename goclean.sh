#!/bin/bash
# The script does automatic checking on a Go package and its sub-packages, including:
# 1. gofmt         (http://golang.org/cmd/gofmt/)
# 2. goimports     (https://github.com/bradfitz/goimports)
# 3. golint        (https://github.com/golang/lint)
# 4. go vet        (http://golang.org/cmd/vet)
# 5. race detector (http://blog.golang.org/race-detector)
# 6. test coverage (http://blog.golang.org/cover)
# ref. https://gist.github.com/hailiang/0f22736320abe6be71ce

set -e

#go get -v golang.org/x/lint/golint
#go get -v golang.org/x/tools/cmd/goimports
#go get -v ./...

# Automatic checks
#test -z "$(gofmt -l -w ./     | tee /dev/stderr)"
#test -z "$(goimports -l -w ./ | tee /dev/stderr)"
#test -z "$(golint ./...       | tee /dev/stderr)"

#go vet ./...
#go test -race ./...

#go test -race ./campaign...
#go test -race ./mails...
#go test -race ./reader...
#go test -race ./utils...

# Run test coverage on each subdirectories and merge the coverage profile.

echo "mode: count" > coverage.cov

ls
# Standard go tooling behavior is to ignore dirs with leading underscors
for dir in $(find . -maxdepth 10 -not -path './.git*' -not -path '*/_*' -not -path './cmd' -not -path './cmd/docs' -not -path './sql' -type d);
do
if ls $dir/*.go &> /dev/null; then
    go test -v -covermode=count -coverprofile=$dir/coverage.tmp $dir
    if [ -f $dir/coverage.tmp ]
    then
        cat $dir/coverage.tmp | tail -n +2 >> coverage.cov
        rm $dir/coverage.tmp
    fi
fi
done

go tool cover -func coverage.cov
cat coverage.cov

# To submit the test coverage result to coveralls.io,
# use goveralls (https://github.com/mattn/goveralls)
# goveralls -covercoverage=coverage.cov -service=travis-ci
