####################################################################################################
# Download dependencies and build
####################################################################################################
FROM golang:1.13.1-alpine3.10 AS dependencies

WORKDIR $GOPATH/src/github.com/caos/boom

# Runtime dependencies
RUN apk update && apk add git curl && \
    mkdir /artifacts && \
    curl -L "https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.4.0/kustomize_v3.4.0_linux_amd64.tar.gz" |tar xvz && \
    mv ./kustomize /artifacts/kustomize && \
    curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.17.0/bin/linux/amd64/kubectl && \
    chmod +x ./kubectl && \
    mv ./kubectl /artifacts/kubectl && \
    curl -L "https://get.helm.sh/helm-v2.12.0-linux-amd64.tar.gz" |tar xvz && \
    mv linux-amd64/helm /artifacts/helm && \
    chmod +x /artifacts/helm

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

# RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o boom main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /boom cmd/boom/*.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /gen cmd/gen-executable/*.go

# ####################################################################################################
# Run binary
# ####################################################################################################
FROM alpine:3.10

RUN apk update && apk add bash ca-certificates
COPY --from=dependencies /artifacts /usr/local/bin/
COPY --from=build /boom /
COPY --from=build /gen /

COPY config/crd /crd
COPY dashboards /dashboards

RUN ./gen

ENTRYPOINT ["/boom"]