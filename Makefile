.PHONY: clean build

SHELL := /bin/bash

NAME=goku
BUILD_DATE=$(shell date '+%Y%m%d')
BUILD_VERSION=$(shell git log -1 --pretty=format:%h)
REGISTRY=ip:8001

test:
	 docker build -f ./docker/Dockerfile.goku -t custom-goku .
	 docker run --name myGoku --link myMysql:myMysql -p 8080:8080 -d custom-goku
	 docker build -f ./docker/Dockerfile.nginx -t custom-nginx .
	 docker run --name myNginx --link myGoku:myGoku -p 80:80 -d custom-nginx

release:
	docker build --no-cache -f ./docker/Dockerfile.goku -t $(REGISTRY)/$(NAME)-server:$(BUILD_DATE)-$(BUILD_VERSION) .
	docker tag $(REGISTRY)/$(NAME)-server:$(BUILD_DATE)-$(BUILD_VERSION) $(REGISTRY)
	docker build --no-cache -f ./docker/Dockerfile.nginx -t $(REGISTRY)/$(NAME)-proxy:$(BUILD_DATE)-$(BUILD_VERSION) .
	docker tag $(REGISTRY)/$(NAME)-proxy:$(BUILD_DATE)-$(BUILD_VERSION) $(REGISTRY)

push:
	docker push $(REGISTRY)/$(NAME)-server:$(BUILD_DATE)-$(BUILD_VERSION)
	docker push $(REGISTRY)/$(NAME)-server
	docker push $(REGISTRY)/$(NAME)-proxy:$(BUILD_DATE)-$(BUILD_VERSION)
	docker push $(REGISTRY)/$(NAME)-proxy