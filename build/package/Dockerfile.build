FROM stakater/go-glide:1.9.3
MAINTAINER "Stakater Team"

RUN apk update

RUN apk -v --update \
    add git build-base && \
    rm -rf /var/cache/apk/* && \
    mkdir -p "$GOPATH/src/github.com/stakater/Forecastle"

ARG PACKR_VERSION=1.11.1
ARG PACKR_FILENAME=packr_${PACKR_VERSION}_linux_386.tar.gz
ARG PACKR_URL=https://github.com/gobuffalo/packr/releases/download/v${PACKR_VERSION}/${PACKR_FILENAME}

RUN mkdir -p /tmp/packr/ && \
    wget ${PACKR_URL} -O /tmp/packr/${PACKR_FILENAME} && \
    tar -xzvf /tmp/packr/${PACKR_FILENAME} -C /tmp/packr/ && \
    mv /tmp/packr/packr /usr/local/bin/packr && \
    rm -rf /tmp/packr

ADD . "$GOPATH/src/github.com/stakater/Forecastle"

RUN cd "$GOPATH/src/github.com/stakater/Forecastle" && \
    glide update && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 packr build -a --installsuffix cgo --ldflags="-s" -o /Forecastle

COPY build/package/Dockerfile.run /

# Running this image produces a tarball suitable to be piped into another
# Docker build command.
CMD tar -cf - -C / Dockerfile.run Forecastle