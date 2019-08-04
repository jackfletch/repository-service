FROM golang:1.12-alpine AS builder
LABEL Author="Jack Fletcher"

# install git, which is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

RUN mkdir -p /go/src/github.com/jackfletch/repository-service/
WORKDIR /go/src/github.com/jackfletch/repository-service/
COPY . .

# fetch dependencies but do not install them [-d]
RUN go get -d -v

# build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' -a -o /go/bin/main

# copy binary to empty image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main /go/bin/main

ENTRYPOINT ["/go/bin/main"]
