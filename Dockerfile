# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial 

ENV p $GOPATH/src/github.com/labbsr0x/bindman-dns-bind9

ADD ./ ${p}
WORKDIR ${p}
RUN go get -v ./...

RUN GIT_COMMIT=$(git rev-parse --short HEAD 2> /dev/null || true) \
 && BUILDTIME=$(TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ') \
 && VERSION=$(git describe --abbrev=0 --tags 2> /dev/null || true) \
 && CGO_ENABLED=0 GOOS=linux go build --ldflags "-s -w \
    -X github.com/labbsr0x/bindman-dns-bind9/version.Version=${VERSION:-unknow-version} \
    -X github.com/labbsr0x/bindman-dns-bind9/version.GitCommit=${GIT_COMMIT} \
    -X github.com/labbsr0x/bindman-dns-bind9/version.BuildTime=${BUILDTIME}" \
    -a -installsuffix cgo -o /bindman-dns-manager src/main.go

# PKG
FROM alpine:latest

RUN apk add --no-cache --update bind-tools

VOLUME [ "/data" ]
COPY --from=builder /bindman-dns-manager /go/bin/

ENTRYPOINT [ "/go/bin/bindman-dns-manager" ]

CMD [ "serve" ]
