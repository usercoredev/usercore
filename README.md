# Usercore (usercore.dev)

### User management system writtefn in Go (GRPC & HTTP) | UNDER DEVELOPMENT

Usercore is under development. It is a user management system written in Go. It provides user registration, login, and
user management features. It uses GRPC and HTTP for communication. It is designed to be used as a microservice.

## Features

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
cp .env.example .env.dev
```

```sh
 docker run --env-file=.env.dev -p 8001:8001 -p 9001:9001 -v $(pwd)/vault:/app/vault usercore/usercore:0.0.1-dev
```

### DB Configuration

> You can use DB_PASSWORD_FILE to load the password from a file instead of setting it in this file directly.
> DB_PASSWORD_FILE will override DB_PASSWORD if both are set.
> Example: DB_PASSWORD_FILE="/run/secrets/db_password"

> DB ENGINE OPTIONS: mysql, postgres, sqlite
> If you want to use sqlite, you should set DB_FILE_PATH.
> DB_FILE_PATH=../development/sqlite.db
> DB_ENGINE=sqlite
