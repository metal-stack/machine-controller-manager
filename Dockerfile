#############      builder                                  #############
FROM golang:1.14.7 AS builder

WORKDIR /go/src/github.com/gardener/machine-controller-manager
COPY . .

RUN .ci/build \
 && strip /go/src/github.com/gardener/machine-controller-manager/bin/rel/machine-controller-manager

#############      base                                     #############
FROM alpine:3.12 as base

RUN apk add --update bash curl tzdata
WORKDIR /

COPY --from=builder /go/src/github.com/gardener/machine-controller-manager/bin/rel/machine-controller-manager /machine-controller-manager
ENTRYPOINT ["/machine-controller-manager"]