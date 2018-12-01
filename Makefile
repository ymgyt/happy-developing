GCP_PROJECT_ID := ${GCP_PROJECT_ID}


init: deps
	circlici setup
	docker login

deps:
	scripts/deps.sh

validate/ci:
	circleci config validate


# docker commands
docker/build:
	docker build -t hpdev:latest -f hpdev/Dockerfile --build-arg GCP_PROJECT_ID=${GCP_PROJECT_ID} hpdev/

docker/run:
	docker run --rm --publish 8123:8123 --name hpdev hpdev:latest

docker/tag:
	docker tag hpdev:latest docker.io/ymgyt/hpdev:latest

docker/push:
	docker push docker.io/ymgyt/hpdev:latest


.phony: init deps validate/ci docker/build docker/run docker/tag docker/push hoge
