#############      builder                                  #############
FROM golang:1.13 AS builder

WORKDIR /go/src/github.com/gardener/machine-controller-manager
COPY . .

RUN .ci/build

#############      machine-controller-manager               #############
FROM alpine:3.10

RUN apk add --update bash curl tzdata
WORKDIR /

COPY --from=builder /go/src/github.com/gardener/machine-controller-manager/bin/rel/machine-controller-manager /machine-controller-manager
ENTRYPOINT ["/machine-controller-manager"]
