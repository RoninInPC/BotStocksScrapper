FROM ubuntu

COPY cacert.pem /etc/ssl/certs/

COPY config.yaml /config/
COPY main /

CMD ["/main"]