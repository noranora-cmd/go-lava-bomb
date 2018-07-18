FROM alpine:3.7
RUN apk add --no-cache openssh ca-certificates

RUN mkdir /app
RUN mkdir /data
WORKDIR /app
COPY go-lava-bomb /app
COPY config.json /app
ENTRYPOINT ["./go-lava-bomb"]
