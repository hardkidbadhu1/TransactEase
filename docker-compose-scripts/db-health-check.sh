#!/bin/sh

# MySQL health check script

# Environment variables
MYSQL_HOST=${MYSQL_HOST}
MYSQL_PORT=${MYSQL_PORT}
MYSQL_USER=${MYSQL_USER}
MYSQL_PASSWORD=${MYSQL_PASSWORD}

# Ping MySQL server
mysqladmin ping -h "$MYSQL_HOST" -P "$MYSQL_PORT" -u "$MYSQL_USER" -p"$MYSQL_PASSWORD"
