BINARY=vault
VERSION=test
CONTAINER_PORT=-p 7007:7014
ENV=-e VAULT_APP_URL=vault.app

build:
	docker run --rm -v ${PWD}:/go/src/app -w /go/src/app lrodham/golang-glide /bin/sh -c "glide install && go build"
	docker build -t ${BINARY}:${VERSION} .
start:
	docker run ${ENV} -d ${CONTAINER_PORT} --name ${BINARY} ${BINARY}:${VERSION}
stop:
	docker stop ${BINARY} && docker rm ${BINARY}
