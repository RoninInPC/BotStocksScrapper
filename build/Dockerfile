FROM ubuntu

COPY cacert.pem /etc/ssl/certs/

RUN apt-get update && apt-get install -y tzdata

ENV TZ="Europe/Moscow"

COPY config.yaml /config/
COPY main /

CMD ["/main"]