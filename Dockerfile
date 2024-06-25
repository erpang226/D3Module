
FROM ubuntu:22.04

RUN mkdir -p /app/config

COPY ./bin/app /app/
COPY ./app.yaml /app/config/
COPY ./bin/app.db /app/

WORKDIR /app

ENTRYPOINT ["/app/app"]
