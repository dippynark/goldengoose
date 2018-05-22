REGISTRY ?= dippynark
APP_NAME ?= goldengoose

docker: docker_build docker_push

docker_build:
	docker build -t $(REGISTRY)/$(APP_NAME) .

docker_push:
	docker push $(REGISTRY)/$(APP_NAME)