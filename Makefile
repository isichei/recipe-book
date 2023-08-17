.PHONY: build deploy build_deploy

build:
	GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o bin/bootstrap main.go

deploy:
	terraform apply -auto-approve

build_deploy: build deploy