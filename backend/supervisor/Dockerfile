FROM golang:alpine AS intermediate

RUN apk update && \
    apk add --no-cache git make

RUN adduser -D -g '' supervisor

WORKDIR $GOPATH/src/

COPY ./backend .
RUN go mod download
RUN go mod verify
RUN cd cmd/supervisor && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -a -o $GOPATH/src/bin/supervisor

FROM scratch

ENV CFG_PATH=/go/bin/config.yml

COPY --from=intermediate /go/src/bin/supervisor /go/bin/supervisor
COPY --from=intermediate /go/src/config.yml /go/bin/config.yml
COPY --from=intermediate /etc/passwd /etc/passwd
COPY --from=intermediate /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER supervisor

WORKDIR /go/bin

EXPOSE 2112

CMD ["/go/bin/supervisor"]