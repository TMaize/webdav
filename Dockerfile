FROM golang:1.20-alpine3.18 as builder

COPY . /app

RUN \
  cd /app && \
  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o webdav . && \
  ls -lh && sleep 1


###############

FROM scratch

EXPOSE 80

WORKDIR /app

COPY --from=builder /app/ /app/

ENTRYPOINT ["/app/webdav"]
