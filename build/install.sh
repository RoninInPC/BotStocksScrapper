#!/bin/bash

rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz && \
export PATH=$PATH:/usr/local/go/bin &&
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . &&
docker build -t bot_docker . &&
docker-compose up -d --build
