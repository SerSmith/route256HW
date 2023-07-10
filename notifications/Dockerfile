FROM ubuntu:22.04

RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates

ADD ./config.yaml /config.yaml
ADD ./bin/app /app


CMD ["/app"]
