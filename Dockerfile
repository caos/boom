####################################################################################################
# Decrypt secrets
####################################################################################################
FROM python:3.7.4-alpine3.10 as secrets

# get secret for decryption from local environment variable
ARG ANSIBLEVAULT_SECRET
RUN echo -n $ANSIBLEVAULT_SECRET >> /secret

# install required libraries for ansible-vault
RUN apk add --no-cache gcc libc-dev libffi-dev openssl-dev
# RUN ln -s /usr/bin/python3 /usr/bin/python
RUN pip3 install --upgrade ansible-vault

RUN mkdir /secrets 
COPY ./build/secretdata/ssh-keys/ /secrets/

# decrypt ssh keys
RUN ansible-vault decrypt --vault-password-file /secret  /secrets/id_rsa-toolsop-build-utils

####################################################################################################
# Download dependencies and build
####################################################################################################
FROM golang:1.13.1-alpine3.10 AS dependencies

# copy all sourcecode from the current repository
RUN mkdir -p $GOPATH/src/github.com/caos/toolsop
COPY ./go.mod $GOPATH/src/github.com/caos/toolsop/go.mod
COPY ./go.sum $GOPATH/src/github.com/caos/toolsop/go.sum
WORKDIR $GOPATH/src/github.com/caos/toolsop

# copy secrets from stage before
RUN mkdir /secrets
COPY --from=secrets /secrets/ /secrets/

# install all tools that are needed for building
RUN apk add --no-cache git openssh-client

# create directory .ssh to add config to, add git public keys to known hosts and copy ssh-keys
RUN mkdir -p $HOME/.ssh && \
    chmod 700 $HOME/.ssh && \
    eval $(ssh-agent) && \
    ssh-keyscan -t rsa github.com >> $HOME/.ssh/known_hosts

# add the mappings of git repository to ssh host with idividual ssh key
# workaround for go get, by which no ssh is used just http
RUN echo -e 'Host github.com-utils\n' \
    '\tHostName github.com\n' \
    '\tIdentityFile /secrets/id_rsa-toolsop-build-utils' >> $HOME/.ssh/config && \ 
    echo -e "[url \"git@github.com-utils:caos/utils\"]\n\tinsteadOf = https://github.com/caos/utils" >> $HOME/.gitconfig 

RUN go env -w GOPRIVATE=github.com/caos/utils,github.com/caos/toolsop && \
    go mod download

# Copy the go source
COPY main.go main.go
COPY api/ $GOPATH/src/github.com/caos/toolsop/api/
COPY controllers/ $GOPATH/src/github.com/caos/toolsop/controllers/
COPY internal/ $GOPATH/src/github.com/caos/toolsop/internal/

# ####################################################################################################
# Run tests
# ####################################################################################################
# FROM dependencies AS test

# WORKDIR $GOPATH/src/github.com/caos/toolsop

# RUN CGO_ENABLED=0 GOOS=linux go test -short $(go list ./... | grep -v /vendor/)
# RUN go test -race -short $(go list ./... | grep -v /vendor/)
# RUN go test -msan -short $(go list ./... | grep -v /vendor/)

# ####################################################################################################
# Run build
# ####################################################################################################
FROM dependencies AS build

WORKDIR $GOPATH/src/github.com/caos/toolsop

# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o toolsop main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o toolsop main.go

# ####################################################################################################
# Run binary
# ####################################################################################################
FROM alpine:3.10 as base
WORKDIR /
COPY --from=build /go/src/github.com/caos/toolsop/toolsop .
COPY tools/kustomize tools/kustomize
COPY tools/toolsets tools/toolsets
COPY tools/start.sh tools/start.sh

RUN apk add curl bash

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

ENTRYPOINT ["/toolsop"]