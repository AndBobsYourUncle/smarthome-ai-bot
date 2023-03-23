FROM ubuntu:22.04

ENV DEBIAN_FRONTEND="noninteractive"

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        golang && \
    update-ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY . /root/smarthome_ai_bot

WORKDIR /root/smarthome_ai_bot
