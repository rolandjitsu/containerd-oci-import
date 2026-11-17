# syntax=docker/dockerfile:1
ARG DEBIAN_VERSION=bookworm-slim
FROM debian:${DEBIAN_VERSION}

CMD ["echo", "HELLO!"]
