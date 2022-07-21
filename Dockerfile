ARG ALPINE_VERSION=3.16
ARG GOLANG_VERSION=1.18

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/test-dexlektika

FROM alpine:${ALPINE_VERSION}

COPY --from=builder /usr/local/bin/test-dexlektika /usr/local/bin/test-dexlektika

ENTRYPOINT ["/usr/local/bin/test-dexlektika"]