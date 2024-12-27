#!/bin/bash

rm -rf /usr/local/go && tar -C /usr/local -xzf go-1.23.0-linux-arm64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin &&
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . &&
docker build -t bot_stock_docker . &&
docker-compose up -d --build
