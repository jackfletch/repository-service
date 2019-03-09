FROM golang:1.12 AS builder
LABEL Author="Jack Fletcher"

# install git, which is required for fetching the dependencies.
RUN apt-get update && apt-get install git ca-certificates && update-ca-certificates

RUN mkdir -p /go/src/github.com/jackfletch/gitRepoService/
WORKDIR /go/src/github.com/jackfletch/gitRepoService/
COPY . .

# fetch dependencies but do not install them [-d]
RUN go get -d -v

# build static binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags '-linkmode external -extldflags "-static"' -a -o /go/bin/main

# copy binary to empty image
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/main /go/bin/main
EXPOSE 8080
ENTRYPOINT ["/go/bin/main"]
