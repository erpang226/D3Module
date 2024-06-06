
FROM ubuntu:22.04

RUN mkdir -p /app/gateway/config

COPY ./bin/gateway /app/gateway/
COPY ./config.yaml /app/gateway/config/
COPY ./bin/app.db /app/gateway/

WORKDIR /app/gateway

ENTRYPOINT ["/app/gateway/gateway"]
