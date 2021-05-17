# weather-api

## Overview
### Get Weather function
Receives the GET request and returns windspeed and temperature

## Setup workspace
### Requirements & Pre-requisites
#### AWS Sam local
To run the API locally, AWS Sam needs to be installed.
https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started.html

#### Postgres
To cache weather data, postgres is needed. This build is tested with postgres 9.6.
Download: https://www.postgresql.org/ftp/source/v9.6.22/
Getting Started: https://www.postgresql.org/docs/9.6/index.html

#### Moq
To generate new mock interfaces, Moq needs to be installed. 
*Note: This isn't required to build or run the project
https://github.com/matryer/moq

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

## Notes:
- Inserting new weather data into database can be done asynchronusly, can invoke a separate lambda function to do it.
- Expanding on above, depending on the amount of traffic expected, we can have separate function dedicated to updating the database on a schedule.
- Tests can be more comprehensive
- Some code can be slimmed down (e.g. test code can be slimmed down via constructor functions for mocks)
- Cache time can be an env variable
- Due to the simplicity of data, this can be done in nosql (i.e. dynamodb) for performance and cost. However, setup overhead is more complex so I just used simple postgres queries
- Similar to above, can use database ORM if database need to be expanded, but no need to over-engineer as of now
- Weatherstack & openweathermap can be more detailed. I didn't spend much time testing out what error codes & responses I can be receiving so the response handler is very generic.
- Didn't spend too much time on implementing gitflow (i.e. develop/release etc branches) or repository configuration
- Didn't spend much time on CI/CD or AWS configurations if you wanted to deploy a live version. But should be easy enough to add if needed
- Can add swagger if api needs to be expanded. Go-swagger allows generating request & response structs (and also validation functions) based on swagger.yaml file.
- Forgot to add city field to DB in v1 and v1.1