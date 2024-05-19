# Use an official Golang runtime as a parent image
FROM golang:1.22.3-alpine3.19

# Install git
RUN apk update && apk add --no-cache git

# Install golang-migrate
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Set the Current Working Directory inside the container
WORKDIR /migrate

# Copy the migration files to the container
COPY db/migrations db/migrations

# change permission of the migration script
RUN chmod +x db/migrations/db-migration.sh




