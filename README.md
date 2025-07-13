# chirpy

boot.dev "Learn HTTP Servers in Go" project  
https://www.boot.dev/lessons/50f37da8-72c0-4860-a7d1-17e4bda5c243

# test

http://localhost:8080/app

# dependencies

## create postgresql DB

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

## gose migration tool

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

`openssl rand -base64 64` to generate the secret

```
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"

PLATFORM="dev"

JWT_SECRET=""
```

# dev dependencies

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
