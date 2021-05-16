# weather-api

## Overview
### Get Weather function
Receives the GET request and returns windspeed and temperature

## Setup workspace
### Dependencies
#### Moq
To generate new mock interfaces, Moq needs to be installed. 
*Note: This isn't required to build or run the project
`https://github.com/matryer/moq`

#### AWS Sam local
To run the API locally, AWS Sam needs to be installed.
`https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-getting-started.html`

### Code & vendor

Git clone this repository.
`git clone git@github.com:TomSED/weather-api.git`

Run `dep ensure` to pull dependencies

### Deployment & Configuration
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