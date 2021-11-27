#!/bin/bash

# https://www.youtube.com/watch?v=G3gnMSyX-XM

docker pull postgres

docker run -p 5433:5432 -d \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=interview \
  postgres

sleep 1s

# create table tokens
psql postgresql://postgres:postgres@localhost:5433/interview -f setup.sql