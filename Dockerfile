#############      builder                                  #############
<<<<<<< HEAD
FROM golang:1.13.4 AS builder
=======
FROM golang:1.13 AS builder
>>>>>>> c72a11c365f072843f594044216f91988125822b

WORKDIR /go/src/github.com/gardener/machine-controller-manager
COPY . .

RUN .ci/build

<<<<<<< HEAD
#############      base                                     #############
FROM alpine:3.10.3 as base
=======
#############      machine-controller-manager               #############
FROM alpine:3.10
>>>>>>> c72a11c365f072843f594044216f91988125822b

RUN apk add --update bash curl tzdata
WORKDIR /

COPY --from=builder /go/src/github.com/gardener/machine-controller-manager/bin/rel/machine-controller-manager /machine-controller-manager
ENTRYPOINT ["/machine-controller-manager"]
