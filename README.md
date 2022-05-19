# AWS Step Functions Test Project

Project to interact with AWS services and test them using [golium](https://github.com/TelefonicaTC2Tech/golium#golium). 
Services tested are:
- Step functions
- Lambda
- S3
- IAM
- Dynamo DB

## Requirements

* [Docker](https://www.docker.com/community-edition)
* [Golang](https://golang.org)
* [AWS CLI](https://aws.amazon.com/es/cli/)
* [SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)

## Quick start
- Install all [requirements](#requirements)
- Install go dependencies
```bash
make install
make download-tools
```
- Launch aws local stack
```bash
docker-compose up
```
- Build aws stack with SAM template
```bash
make build
```
- Zip deployment
```bash
make zip-deployment
```
- Launch test
```bash
make test
```

## Other commands
- Validate SAM template
```bash
make validate-template
```
- Lint go code
```bash
make lint
```

- Clean zip Deployment
```bash
make clean-zip-deployment
```
