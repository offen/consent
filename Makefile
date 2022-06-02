build:
	@mkdir -p bin
	@go build -o bin/consent-$$(uname -s | tr '[:upper:]' '[:lower:]')-$$(uname -p) cmd/consent/main.go

DOCKER_TAG ?= local
docker:
	@docker build \
		-t offen/consent:$(DOCKER_TAG) .

up:
	@docker-compose up
