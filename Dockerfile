####################################################################################################
# Download dependencies and build
####################################################################################################
FROM golang:1.13.1-alpine3.10 AS dependencies

ENV GO111MODULE on
WORKDIR $GOPATH/src/github.com/caos/toolsop

# copy all sourcecode from the current repository
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy the go source
COPY cmd cmd
COPY api api
COPY controllers controllers
COPY internal internal

# ####################################################################################################
# Run tests
# ####################################################################################################
# FROM dependencies AS test

# RUN CGO_ENABLED=0 GOOS=linux go test -short $(go list ./... | grep -v /vendor/)
# RUN go test -race -short $(go list ./... | grep -v /vendor/)
# RUN go test -msan -short $(go list ./... | grep -v /vendor/)

# ####################################################################################################
# Run build
# ####################################################################################################
FROM dependencies AS build

# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o toolsop main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /toolsop cmd/toolsop/*.go

# ####################################################################################################
# Run binary
# ####################################################################################################
FROM alpine:3.10
WORKDIR /
RUN apk update && apk add bash curl

RUN curl -L "https://get.helm.sh/helm-v2.12.0-linux-amd64.tar.gz" |tar xvz && \
    mv linux-amd64/helm /usr/bin/helm && \
    chmod +x /usr/bin/helm && \
    rm -rf linux-amd64 && \
    rm -f /var/cache/apk/* 

RUN curl -s "https://api.github.com/repos/kubernetes-sigs/kustomize/releases" |\
    grep browser_download |\
    grep linux |\
    cut -d '"' -f 4 |\
    grep /kustomize/v |\
    sort | tail -n 1 |\
    xargs curl -O -L && \
    tar xzf ./kustomize_v*_linux_amd64.tar.gz && \
    mv kustomize /usr/local/bin/kustomize

RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /usr/local/bin/kubectl

RUN apk del curl

COPY --from=build /toolsop /

COPY tools/kustomize tools/kustomize
COPY tools/toolsets tools/toolsets
COPY tools/start.sh tools/start.sh

ENTRYPOINT ["/toolsop"]