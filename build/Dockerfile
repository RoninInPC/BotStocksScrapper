FROM scratch

COPY cacert.pem /etc/ssl/certs/

COPY ../config/config.yaml /config/
COPY main /

CMD ["/main"]