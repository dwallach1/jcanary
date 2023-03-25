
build:
	docker buildx build --platform=linux/amd64 -t davidwallach/jcanary:latest .

run:
	docker run -it --rm --name jcanary -e RULES_CONFIG="./rules2.json" davidwallach/jcanary

push:
	docker logout
	docker login -u="${DOCKER_USER}" -p="${DOCKER_PASSWORD}" docker.io
	docker push davidwallach/jcanary:latest


run-test-web-server:
	docker build -t jcanary-web-server -f ./example/Dockerfile ./example
	docker run --rm --name web-server -p 8080:8080 jcanary-web-server