FROM golang:1.24.2-alpine



COPY serve.sh /usr/local/bin/serve
RUN chmod +x /usr/local/bin/serve

COPY code-runner.sh /usr/local/bin/code-runner
RUN chmod +x /usr/local/bin/code-runner

COPY create-admin.sh /usr/local/bin/create-admin
RUN chmod +x /usr/local/bin/create-admin

COPY load-test-data.sh /usr/local/bin/load-test-data
RUN chmod +x /usr/local/bin/load-test-data

ADD ./src /src
WORKDIR /src
RUN go mod download
RUN go build -o judge ./
