FROM golang:1.17.1-alpine3.13 AS base

ENV GOPROXY=direct \
    GOSUMDB=off \
    GOPATH=/opt/service/.go

RUN apk add --no-cache \
            make==4.3-r0 \
            gcc==10.2.1_pre1-r3 \
            libc-dev==0.7.2-r3 \
            git==2.30.2-r0
WORKDIR /opt/service/

COPY . .
COPY Makefile .

RUN make build

FROM scratch
USER nobody
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /opt/service/bin /bin

HEALTHCHECK --interval=1s --timeout=1s --start-period=2s --retries=3 CMD [ "/bin/healthcheck" ]
