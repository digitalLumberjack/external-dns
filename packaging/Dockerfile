FROM alpine:3.2
MAINTAINER Rancher Labs, Inc.
RUN apk add --update ca-certificates 

ENV EXT_DNS_RELEASE v0.4.1
ADD external-dns /usr/bin/external-dns

ENTRYPOINT ["/usr/bin/external-dns"]
