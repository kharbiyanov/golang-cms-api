FROM alpine:latest
RUN apk add -U --no-cache ca-certificates
WORKDIR /app
VOLUME /app
CMD ["main"]