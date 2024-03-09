# Usercore (usercore.dev)

### User management system writtefn in Go (GRPC & HTTP) | UNDER DEVELOPMENT

Usercore is under development. It is a user management system written in Go. It provides user registration, login, and
user management features. It uses GRPC and HTTP for communication. It is designed to be used as a microservice.

## Features

- POSTGRES, MYSQL, SQLITE support
- GRPC & HTTP support
- JWT based authentication
- Email & SMS verification

## Endpoints
- Sign in
- Sign up
- Sign out
- Password reset
- Password reset confirmation
- Refresh token
- Verify token
- Get user
- Update user
- Delete user
- List users
- Change password
- Change email
- Send verification code (Email & SMS)
- Verify code (Email & SMS)
- Get sessions
- Revoke session

## TODOs

- [ ] Add tests
- [ ] Add more documentation
- [ ] Add more examples
- [ ] Role & permission management
- [ ] Social login (Google, Facebook, Twitter, etc)

## How to run

```sh
cp .env.example .env
```

```sh
 docker run --env-file=.env -p 8001:8001 -p 9001:9001 -v $(pwd)/vault:/app/vault usercore/usercore:0.0.1-dev
```

### DB Configuration

> You can use DB_PASSWORD_FILE to load the password from a file instead of setting it in this file directly.
> DB_PASSWORD_FILE will override DB_PASSWORD if both are set.
> Example: DB_PASSWORD_FILE="/run/secrets/db_password"

> DB ENGINE OPTIONS: mysql, postgres, sqlite
> If you want to use sqlite, you should set DB_FILE_PATH.
> DB_FILE_PATH=../development/sqlite.db
> DB_ENGINE=sqlite


### Example docker-compose.yaml

Note: If you want to build the image locally, you can use the `build.sh` script.

```sh
./build.sh -t <version> # Example: ./build.sh -t 0.0.1-dev
```

After building the image, you can use the following one of the examples `<db_engine>-docker-compose.yaml` file to run the service.
If you build the image locally, you should change the image name to docker image hash.

NOTE: Don't forget to set the correct environment variables in the `.env` file.

**[mysql-docker-compose.yaml](./examples/mysql-docker-compose.yaml)**

**[postgres-docker-compose.yaml](./examples/postgres-docker-compose.yaml)**

**[sqlite-docker-compose.yaml](./examples/sqlite-docker-compose.yaml)**
