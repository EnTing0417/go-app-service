IMAGE_NAME := go-app-service
REPO_NAME := enting0417/go-app-service
TAG_NAME := v1

.PHONY: build

build:
	sudo docker build -t $(IMAGE_NAME) .
	sudo docker tag $(IMAGE_NAME) $(REPO_NAME):$(TAG_NAME)
	sudo docker login
	sudo docker push $(REPO_NAME):$(TAG_NAME)