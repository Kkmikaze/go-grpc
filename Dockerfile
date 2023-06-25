FROM golang:1.19-buster as builder

COPY . /root/go/src/app

ARG BUILD_VERSION

WORKDIR /root/go/src/app

ENV BUILD_VERSION=$BUILD_VERSION

RUN go mod vendor

RUN CGO_ENABLED=0 GOOS=linux go build -a -v -ldflags "-X main.version=$(BUILD_VERSION)" -installsuffix cgo -o app ./cmd/server

FROM bitnami/minideb:buster

WORKDIR /usr/app

COPY --from=builder /root/go/src/app/app /usr/app/app

ENTRYPOINT ["/usr/app/app"]
CMD ["service", "--authport", "4010", "--movieport", "4011", "--gwport", "4000"]