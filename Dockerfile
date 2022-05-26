# Copyright 2022 - Offen Authors <hioffen@posteo.de>
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.18-alpine as builder

ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /code
COPY go.mod go.sum /code/
RUN go mod download
COPY . /code/
RUN GOARCH=$TARGETARCH${TARGETVARIANT:+-$TARGETVARIANT} go build -o consent cmd/consent/main.go

FROM alpine:3.15
LABEL maintainer="offen <hioffen@posteo.de>"

ARG TARGETARCH
ARG TARGETVARIANT

RUN addgroup -g 10001 -S consent \
	&& adduser -u 10000 -S -G consent -h /home/consent consent
RUN apk add -U --no-cache ca-certificates libcap tini bind-tools

COPY --from=builder ./code/consent /opt/consent/consent-linux-$TARGETARCH${TARGETVARIANT:+-$TARGETVARIANT}
RUN setcap CAP_NET_BIND_SERVICE=+eip /opt/consent/consent-linux-$TARGETARCH${TARGETVARIANT:+-$TARGETVARIANT}
RUN ln -s /opt/consent/consent-linux-$TARGETARCH${TARGETVARIANT:+-$TARGETVARIANT} /sbin/consent

ENV PORT 80
EXPOSE 80 443

HEALTHCHECK --interval=1m --timeout=5s \
  CMD wget -qO- http://localhost:80/healthz || exit 1

ENTRYPOINT ["/sbin/tini", "--", "consent"]

USER consent
WORKDIR /home/consent
