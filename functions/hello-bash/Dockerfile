FROM alpine
WORKDIR /app

RUN apk add bash jq
ADD entrypoint.sh /app

ENTRYPOINT ["bash", "/app/entrypoint.sh"]
