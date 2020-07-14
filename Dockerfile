FROM alpine:latest
RUN apk add curl
COPY ./bin/linux/amd64/k0s /usr/bin/k0s
ENTRYPOINT ["/usr/bin/k0s"]
