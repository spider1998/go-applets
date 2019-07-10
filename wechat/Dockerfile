FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime

COPY binaries/wechat /usr/bin

WORKDIR /app

EXPOSE 80

CMD ["wechat"]
