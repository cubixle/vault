FROM golang:1 AS build-env
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 go build -o /app/vault /app/main.go

FROM alpine:latest
COPY --from=build-env /app/vault /vault
RUN ls -lah /vault
CMD ["/vault"]
