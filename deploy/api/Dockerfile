# Step 1:
FROM golang:1.18.2-alpine AS builder

RUN apk update && apk add --no-cache git make

WORKDIR $GOPATH/src/github.com/iris-contrib/go-rel-iris-example
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64\
    go build -mod=vendor -ldflags="-w -s" -o /go/bin/api ./cmd/api

# Step 2:
# you can also use scratch here, but I prefer to use alpine because it comes with basic command such as curl useful for debugging.
FROM alpine

RUN apk update && apk add --no-cache curl ca-certificates
RUN rm -rf /var/cache/apk/*

COPY --from=builder --chown=65534:0 /go/bin/api /go/bin/api

USER 65534
EXPOSE 3000

ENTRYPOINT ["/go/bin/api"]
