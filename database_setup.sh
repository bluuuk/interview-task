# https://www.youtube.com/watch?v=G3gnMSyX-XM

docker pull postgres

docker run -p 5432:5432 -d \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=interview \
  postgres

psql interview -h localhost -U postgres

CREATE TABLE TOKENS(
  token VARCAHR(7) NOT NULL
)