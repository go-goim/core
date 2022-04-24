FROM alpine:latest

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g" /etc/apk/repositories
RUN apk add --no-cache vim
RUN apk add --no-cache net-tools
RUN apk add --no-cache procps
RUN apk add --no-cache curl
RUN apk add --no-cache iputils

WORKDIR /app

COPY bin/msg .
COPY apps/msg/config config

CMD ["./msg","--conf", "./config"]