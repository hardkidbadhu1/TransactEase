FROM golang:1.22.3-alpine3.19

RUN apk update && apk add --no-cache git bash postgresql-client

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /migrate

COPY db/migrations /migrate/db/migrations
COPY docker-compose-scripts/wait-for-db.sh /migrate/wait-for-db.sh
COPY db/migrations/db-migration.sh /migrate/db-migration.sh

RUN chmod +x /migrate/wait-for-db.sh
RUN chmod +x /migrate/db-migration.sh
