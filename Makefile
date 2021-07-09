DOCKER_IMAGE = dippynark/goldengoose:$(shell git branch --show-current)

docker_build:
	docker build -t $(DOCKER_IMAGE) .

docker_push: docker_build
	docker push $(DOCKER_IMAGE)
