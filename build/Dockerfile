FROM scratch

COPY cacert.pem /etc/ssl/certs/

COPY config.ini /config/
COPY main /

CMD ["/main"]