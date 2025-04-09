FROM golang:tip-alpine3.21
ADD ./src /src
WORKDIR /src

COPY serve.sh /usr/local/bin/serve
RUN chmod +x /usr/local/bin/serve

COPY code-runner.sh /usr/local/bin/code-runner
RUN chmod +x /usr/local/bin/code-runner