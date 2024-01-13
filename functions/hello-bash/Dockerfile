FROM alpine
WORKDIR /app

RUN apk add jq
ADD entrypoint.sh /app

ENTRYPOINT ["sh", "/app/entrypoint.sh"]
