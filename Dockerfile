# A minimal Docker image based on Alpine Linux with a complete package index and only 5 MB in size!
FROM alpine

# expose port
EXPOSE 80

# Add executable into image
COPY build/camelcaseapp /
COPY credentials.json /

RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true


CMD ["/camelcaseapp"]

