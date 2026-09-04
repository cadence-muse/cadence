FROM alpine:3.20

RUN apk add --no-cache ca-certificates curl && \
    update-ca-certificates && \
    addgroup -g 1001 cadenceuser && \
    adduser -u 1001 -D -G cadenceuser -s /sbin/nologin -g "go service user" cadenceuser

COPY ./bin/cadence /app/bin/
WORKDIR /app

EXPOSE 8080

USER cadenceuser

ENTRYPOINT [ "/app/bin/cadence" ]
CMD ["service"]
