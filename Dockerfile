FROM ubuntu:22.04

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        golang && \
    apt-get clean && \
    update-ca-certificates && \
    rm -rf /var/lib/apt/lists/*
