FROM golang:1.22.3-alpine3.19

LABEL maintainer="Badhmanaban M <badhmanabanm@hotmail.com>"

RUN apk update && apk add --no-cache git bash postgresql-client

WORKDIR /app

COPY docker-compose-scripts/wait-for-db.sh /app/wait-for-db.sh
COPY configuration/config.json configuration/config.json

COPY bin/transact-api transact-api

RUN chmod +x /app/wait-for-db.sh
RUN chmod +x /app/transact-api

EXPOSE 8080

CMD ["./transact-api"]