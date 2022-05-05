FROM alpine:edge

# RUN \
#  echo "http://dl-cdn.alpinelinux.org/alpine/latest-stable/main" >> /etc/apk/repositories \
#  echo "http://dl-cdn.alpinelinux.org/alpine/latest-stable/community" >> /etc/apk/repositories

RUN apk update && \
 apk add --no-cache \
    mongodb \
    mongodb-tools \
    bash

# FROM mongodb-server

WORKDIR /
ENTRYPOINT [ "/bin/bash" ]