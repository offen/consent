build:
	@mkdir -p bin
	@go build -o bin/consent-linux-amd64 cmd/consent/main.go

DOCKER_TAG ?= local
DOCKER_TARGETARCH ?= amd64
DOCKER_TARGETVARIANT ?=
docker:
	@docker build \
		--build-arg=TARGETARCH=$(DOCKER_TARGETARCH) \
		--build-arg=TARGETVARIANT=$(DOCKER_TARGETVARIANT) \
		-t offen/consent:$(DOCKER_TAG) .

up:
	@docker-compose up
