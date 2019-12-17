export CGO_ENABLED:=0

VERSION=$(shell ./scripts/git-version.bash)
DOCKER_REPO=hetznercloud/hcloud-cloud-controller-manager

all: build

build: clean bin/hcloud-cloud-controller-manager

bin/%:
	@go build -o bin/$* .

container: build
	docker build -t docker.be-mobile.biz:5000/hcloud-cloud-controller-manager:v1.2.0-test1 .

release-container:
	docker push docker.be-mobile.biz:5000/hcloud-cloud-controller-manager:v1.2.0-test1

test:
	@./scripts/test.bash

clean:
	@rm -rf bin/*

.PHONY: all build clean test container release-container
