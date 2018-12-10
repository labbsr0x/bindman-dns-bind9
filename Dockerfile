# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial 

ENV BUILD_PATH=$GOPATH/src/github.com/labbsr0x/bindman-dns-bind9/src

RUN mkdir -p ${BUILD_PATH}
WORKDIR ${BUILD_PATH}

ADD ./src ./
RUN go get -v ./...

WORKDIR ${BUILD_PATH}/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /manager .

# PKG
FROM alpine:latest

RUN apk add --no-cache --update \
  curl \
  wget \
  nmap \
  bind-tools

COPY --from=builder /manager /

VOLUME [ "/data" ]

EXPOSE 7070

ENV BINDMAN_NAMESERVER_ADDRESS ""
ENV BINDMAN_NAMESERVER_PORT ""
ENV BINDMAN_NAMESERVER_KEYFILE ""
ENV BINDMAN_NAMESERVER_ZONE ""
ENV BINDMAN_MODE ""

CMD ["./manager"]
