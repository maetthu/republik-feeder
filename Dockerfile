FROM alpine:latest
LABEL org.opencontainers.image.source=https://github.com/maetthu/republik-feeder
LABEL org.opencontainers.image.description="republik-feeder"
LABEL org.opencontainers.image.licenses=MIT

RUN apk add --no-cache ca-certificates

COPY /republik-feeder /usr/bin/republik-feeder

EXPOSE 8080/tcp
ENTRYPOINT ["/usr/bin/republik-feeder", ":8080"]
CMD []
