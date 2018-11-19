# BUILD
FROM golang:1.11-alpine as builder

RUN apk add --no-cache git mercurial 

ENV BUILD_PATH=$GOPATH/src/github.com/labbsr0x/sandman-bind9-manager/src

RUN mkdir -p ${BUILD_PATH}
WORKDIR ${BUILD_PATH}

ADD ./src ./
RUN go get -v ./...

WORKDIR ${BUILD_PATH}/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /manager .

# PKG
FROM ubuntu:bionic-20180526 AS add-apt-repositories

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y gnupg \
 && apt-key adv --fetch-keys http://www.webmin.com/jcameron-key.asc \
 && echo "deb http://download.webmin.com/download/repository sarge contrib" >> /etc/apt/sources.list

# RUNTIME OS - Manager should be run alongside nsupdate
FROM ubuntu:bionic-20180526

ENV BIND_VERSION=9.11.3

COPY --from=add-apt-repositories /etc/apt/trusted.gpg /etc/apt/trusted.gpg

COPY --from=add-apt-repositories /etc/apt/sources.list /etc/apt/sources.list

COPY --from=builder /manager /

RUN rm -rf /etc/apt/apt.conf.d/docker-gzip-indexes \
 && apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y dnsutils \
 && rm -rf /var/lib/apt/lists/*

EXPOSE 7070

CMD ["./manager"]
