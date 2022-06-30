FROM node:16 as build-deps
WORKDIR /usr/src/app
COPY ./frontend ./
RUN yarn install --network-timeout 1000000 && yarn build

# Build the manager binary
FROM --platform=${BUILDPLATFORM} golang:1.18 as builder

ARG TARGETOS
ARG TARGETARCH

RUN mkdir -p "$GOPATH/src/github.com/stakater/Forecastle"

WORKDIR "$GOPATH/src/github.com/stakater/Forecastle"

# Copy manifests
COPY . .

# Install Packr2
ARG PACKR_VERSION=2.7.1
ARG PACKR_FILENAME=packr_${PACKR_VERSION}_linux_386.tar.gz
ARG PACKR_URL=https://github.com/gobuffalo/packr/releases/download/v${PACKR_VERSION}/${PACKR_FILENAME}

RUN mkdir -p /tmp/packr/ && \
    wget ${PACKR_URL} -O /tmp/packr/${PACKR_FILENAME} && \
    tar -xzvf /tmp/packr/${PACKR_FILENAME} -C /tmp/packr/ && \
    mv /tmp/packr/packr2 /usr/local/bin/packr2 && \
    rm -rf /tmp/packr

# Copy dependencies
COPY --from=build-deps /usr/src/app/build ./frontend/build/

# Build
RUN go mod download && \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GO111MODULE=on packr2 build -a --installsuffix cgo --ldflags="-s" -o /Forecastle && \
    packr2 clean

# Use distroless as minimal base image to package the Forecastle binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /Forecastle .
USER nonroot:nonroot

ENTRYPOINT ["/Forecastle"]
