FROM golang:alpine as builder

ADD ./main.go /go/src/github.com/cirocosta/l4/main.go
ADD ./lib /go/src/github.com/cirocosta/l4/lib
ADD ./vendor /go/src/github.com/cirocosta/l4/vendor

WORKDIR /go/src/github.com/cirocosta/l4

RUN set -ex && \
  CGO_ENABLED=0 go build -tags netgo -v -a -ldflags '-extldflags "-static"' && \
  mv ./l4 /usr/bin/l4

FROM busybox
COPY --from=builder /usr/bin/l4 /usr/local/bin/l4

ENTRYPOINT [ "l4" ]


