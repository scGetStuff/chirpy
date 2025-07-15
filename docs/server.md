# run time dependencies

you need to create the DB before running the server

## install and create postgresql DB

```shell
sudo apt update
sudo apt install postgresql postgresql-contrib
psql --version

# OS user
sudo passwd postgres
sudo service postgresql start

sudo -u postgres psql
CREATE DATABASE chirpy;
\c chirpy
# DB user
ALTER USER postgres PASSWORD 'postgres';
SELECT version();
exit
```

## goose migration tool to create tables

install

```shell
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -version
```

run

```shell
goose postgres postgres://postgres:postgres@localhost:5432/chirpy down
goose postgres postgres://postgres:postgres@localhost:5432/chirpy up
```

## create `.env` file

```
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"

# dev enables the reset endpoint
PLATFORM="dev"

# `openssl rand -base64 64` to generate the secret
JWT_SECRET=""

# this was defined in the lesson, I don't want to show it in documentation
POLKA_KEY=""
```

# dev dependencies

to generate the `internal/database` package

## sqlc

install

```shell
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
sqlc version
```

run

```shell
sqlc generate
```
