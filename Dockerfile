FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y --no-install-suggests --no-install-recommends ca-certificates curl && \
    apt-get clean && \
    groupadd -g 1001 cadenceuser && \
    useradd -u 1001 -r -g 1001 -s /sbin/nologin -c "go service user" cadenceuser

RUN update-ca-certificates --fresh

ADD ./bin/cadence /app/bin/
WORKDIR /app

EXPOSE 8080

USER cadenceuser

ENTRYPOINT [ "/app/bin/cadence" ]
CMD ["service"]
