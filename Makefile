GIT_COMMIT=$(shell git rev-list -1 HEAD)
LDFLAGS=-a -tags -ldflags="-w extldflags '-static' -X lesson8/request.GitCommit=${GIT_COMMIT}" -o ./app ./
FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

.PHONY: test
test:
	go test ./...

.PHONY: build
build: test
	${FLAGS} go build ${LDFLAGS}
