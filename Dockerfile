FROM golang:alpine as builder
ENV GO111MODULE on
COPY . $GOPATH/src/github.com/sourcegraph/shortcut
WORKDIR $GOPATH/src/github.com/sourcegraph/shortcut
RUN apk add --no-cache ca-certificates
RUN apk add --no-cache git && go get -d . && apk del git
RUN env GOOS=linux GARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/shortcut .
RUN adduser -D -g '' appuser

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/shortcut /go/bin/shortcut
USER appuser
EXPOSE 3980
ENTRYPOINT ["/go/bin/shortcut"]