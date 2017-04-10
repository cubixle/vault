BINARY=vault
VERSION=test
CONTAINER_PORT=-p 7007:7014
ENV=-e VAULT_APP_URL=vault.app

build:
	docker run --rm -v ${PWD}:/go/src/app -w /go/src/app -e GOOS=linux -e GOARCH=386 sipsynergy/go-builder /bin/sh -c "godep get && godep go build"
	docker build -t ${BINARY}:${VERSION} .
start:
	docker run -e ${ENV} -d ${CONTAINER_PORT} --name ${BINARY} ${BINARY}:${VERSION}
stop:
	docker stop ${BINARY} && docker rm ${BINARY}
