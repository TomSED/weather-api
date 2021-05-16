S3_BUCKET := $(S3_BUCKET)
STACK_NAME := $(STAGE)-WEATHERAPI
TEMPLATE := template.yaml
PACKAGED_TEMPLATE := dist/$(STACK_NAME)-template.yaml
WORKERS := $(addprefix dist/,$(notdir $(wildcard workers/*)))
VARS := Stage=$(STAGE) WeatherStackApiKey=$(WEATHERSTACK_API_KEY) OpenWeatherMapApiKey=$(OPENWEATHERMAP_API_KEY)

.PHONY: clean deps

clean:
	rm -rf ./dist

deps:
	go mod vendor

build: clean deps $(WORKERS)

$(WORKERS): vendor
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ $(addprefix ./workers/,$(notdir $@))

test:
	go test $(shell go list ./... ) -coverprofile c.out

local: clean build
	$(VARS) sam local start-api -p 8080

local-faster: clean build
	$(VARS) sam local start-api --skip-pull-image -p 8080

$(PACKAGED_TEMPLATE): build
	aws cloudformation package --template-file $(TEMPLATE) --s3-bucket $(S3_BUCKET) --output-template-file $(PACKAGED_TEMPLATE)

deploy: $(PACKAGED_TEMPLATE)
	aws cloudformation deploy --stack-name $(STACK_NAME) \
	--template-file $(PACKAGED_TEMPLATE) \
	--capabilities CAPABILITY_IAM \
	--parameter-override $(VARS)
	aws cloudformation describe-stacks \
    --stack-name $(STACK_NAME) \
    --query 'Stacks[].Outputs'
