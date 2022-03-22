BINARY=vault
VERSION=test
VAULT_APP_URL="vault.app"
HOST_PORT="8080"
VAULT_PORT="8080"

build:
	docker build -t ${BINARY}:${VERSION} .
start:
	docker run -e VAULT_APP_URL=${APP_URL} -e VAULT_PORT=${VAULT_PORT} -d -p ${HOST_PORT}:${VAULT_PORT} --name ${BINARY} ${BINARY}:${VERSION}
stop:
	docker stop ${BINARY} && docker rm ${BINARY}
