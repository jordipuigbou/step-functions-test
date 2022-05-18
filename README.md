# AWS Step Functions Test Project

Project to interact with AWS services and test them using [golium](https://github.com/TelefonicaTC2Tech/golium#golium). 
Services tested are:
- Step functions
- Lambda
- S3
- IAM
- Dynamo DB

## Quick start

- Launch aws local stack
```bash
docker-compose up
```

- Build with SAM
```bash
make build
```

- Zip Deployment
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

- Clean zip Deployment
```bash
make clean-zip-deployment
```

## Pending improvements
- Update aws golang sdk to v2