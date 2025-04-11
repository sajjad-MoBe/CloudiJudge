FROM golang:1.24.2-alpine



COPY serve.sh /usr/local/bin/serve
RUN chmod +x /usr/local/bin/serve

COPY code-runner.sh /usr/local/bin/code-runner
RUN chmod +x /usr/local/bin/code-runner


ADD ./src /src
WORKDIR /src
RUN go mod download
