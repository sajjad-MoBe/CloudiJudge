FROM golang:1.24.2-alpine

RUN apk add --no-cache bash

WORKDIR /app

COPY runner.sh .

RUN chmod +x runner.sh

RUN adduser -D -u 1000 appuser
USER appuser

ENTRYPOINT ["./runner.sh"]