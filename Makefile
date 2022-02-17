LDFLAGS=-a -tags netgo -ldflags='-w -extldflags "-static"' -o ./app ./server
FLAGS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

.PHONY: test
test:
	go test ./...

.PHONY: build
build: test
	${FLAGS} go build ${LDFLAGS}
