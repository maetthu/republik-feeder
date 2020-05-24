# Build image
FROM alpine:latest AS build

# Build requirements
RUN apk add --no-cache ca-certificates

# Copy binary
COPY republik-feeder /republik-feeder

# ---

# Runtime image
FROM scratch
LABEL maintainer="Matthias Blaser <git@mooch.ch>"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /republik-feeder /republik-feeder

EXPOSE 8080/tcp
ENTRYPOINT ["/republik-feeder", ":8080"]
CMD []
