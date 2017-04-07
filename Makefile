BINARY=vault
VERSION=test
CONTAINER_PORT=-p 7007:7014

build:
	docker run --rm -v ${PWD}:/go/src/app -w /go/src/app -e GOOS=linux -e GOARCH=386 sipsynergy/go-builder /bin/sh -c "godep get && godep go build"
	docker build -t ${BINARY}:${VERSION} .
start:
	docker run -d ${CONTAINER_PORT} --name ${BINARY} ${BINARY}:${VERSION}
stop:
	docker stop ${BINARY} && docker rm ${BINARY}
