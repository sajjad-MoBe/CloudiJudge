FROM golang:1.24.2-alpine

RUN apk add --no-cache bash

WORKDIR /app

COPY builder.sh /usr/local/bin/builder
RUN chmod +x /usr/local/bin/builder


COPY runner.sh /usr/local/bin/runner
RUN chmod +x /usr/local/bin/runner


RUN adduser -D -u 1000 runner
USER runner

ENTRYPOINT ["builder"]