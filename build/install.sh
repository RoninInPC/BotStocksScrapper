#!/bin/bash

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin && \
cd .. && cd cmd && \
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . && \
cd .. && cp cmd/main build && cp config/config.yaml build &&\
docker build -t bot_stocks_docker . && \
docker-compose up -d --build
