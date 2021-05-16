# weather-api

## Overview
### Get Weather function
Receives the GET request and returns windspeed and temperature

## Setup workspace
### Requirements & Pre-requisites
#### AWS Sam local
To run the API locally, AWS Sam needs to be installed.
`https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started.html`

#### Postgres
To cache weather data, postgres is needed. This build is tested with postgres 9.6.
Download: `https://www.postgresql.org/ftp/source/v9.6.22/`
Getting Started: `https://www.postgresql.org/docs/9.6/index.html`

#### Moq
To generate new mock interfaces, Moq needs to be installed. 
*Note: This isn't required to build or run the project
`https://github.com/matryer/moq`

### Code & vendor
Git clone this repository.
`git clone git@github.com:TomSED/weather-api.git`

Run `go mod vendor` to pull dependencies

### Deployment & Configuration
#### Setup Postgres
The current set up script creates a 'weather' table, make sure you don't have a conflicting table name in your db.
1. Create a `/.env` file according to `/.env.template`.
```bash
$ export $(grep -v '^#' .env | xargs)
```

2. Setup DB
```bash
$ go run pkg/postgres/setup/main.go
```

#### Deploy to AWS
1. Create a `/.env` file according to `/.env.template`.
```bash
$ export $(grep -v '^#' .env | xargs)
```

2. Deploy
```bash
$ make deploy
```

#### Deploy to localhost:8080
1. Create a `/.env` file according to `/.env.template`.
```bash
$ export $(grep -v '^#' .env | xargs)
```

2. Deploy
```bash
$ make local
```