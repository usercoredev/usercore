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


### Example docker-compose.yaml

Note: If you want to build the image locally, you can use the `build.sh` script.

```sh
./build.sh -t <version> # Example: ./build.sh -t 0.0.1-dev
```

After building the image, you can use the following `docker-compose.yaml` file to run the service.
If you build the image locally, you should change the image name to docker image hash.

```yaml
version: "3.8"

services:
  usercore_app:
    image: usercore/usercore:0.0.1-dev # Change this to the image hash if you build the image locally
    container_name: usercore_app
    depends_on:
      usercore_database:
        condition: service_healthy
    env_file:
      - ./service/.env.dev
    networks:
      - usercore_net
    #volumes:
    # - app_data:/app
    ports:
      - 8000:8000
      - 9000:9000
    secrets:
      - jwt_private_key
      - jwt_public_key
      - clients
      - db_password
      - server_cert.pem
      - cache_password
    deploy:
      restart_policy:
        condition: on-failure
  usercore_cache:
    image: redis:latest
    container_name: usercore_cache
    hostname: usercore_cache
    networks:
      - usercore_net
    deploy:
      restart_policy:
        condition: on-failure
    secrets:
      - cache_password
    environment:
      REDIS_PASSWORD: run/secrets/cache_password
      REDIS_USER: usercore
      REDIS_PORT: 6379
    #volumes:
    #  - cache_data:/data
    ports:
      - 6379:6379
  usercore_database:
    image: mariadb:latest
    container_name: usercore_database
    hostname: usercore_db
    networks:
      - usercore_net
    healthcheck:
      interval: 5s
      retries: 10
      test:
        [
          "CMD",
          "healthcheck.sh",
          "--su-mysql",
          "--connect",
          "--innodb_initialized",
        ]
      timeout: 30s
    secrets:
      - db_password
      - source: server_cert.pem
        target: /etc/mysql/certs/server-cert.pem
      - source: server_key.pem
        target: /etc/mysql/certs/server-key.pem
    #volumes:
    #  - db_data:/var/lib/mysql
    command: >
      --ssl-cert=/etc/mysql/certs/server-cert.pem
      --ssl-key=/etc/mysql/certs/server-key.pem
    environment:
      MARIADB_USER: usercore
      MARIADB_DATABASE: usercore
      MARIADB_PASSWORD_FILE: /run/secrets/db_password
      MARIADB_PORT: 3306
      MARIADB_RANDOM_ROOT_PASSWORD: true
      MARIADB_CHARSET: utf8
secrets:
  jwt_private_key:
    file: ./service/vault/example/jwt.private
  jwt_public_key:
    file: ./service/vault/example/jwt.public
  clients:
    file: ./service/vault/example/clients.json
  server_cert.pem:
    file: ./service/vault/example/certs/server-cert.pem
  server_key.pem:
    file: ./service/vault/example/certs/server-key.pem
  cache_password:
    file: ./service/vault/example/redis-pass.txt
  db_password:
    file: ./service/vault/example/db-pass.txt
networks:
  usercore_net:
    driver: bridge
#volumes:
#db_data:
#cache_data:
#app_data:

```
