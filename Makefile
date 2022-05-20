.PHONY: generate
generate:
	@docker-compose run --rm server go generate

.PHONY: up
up: generate
	@docker-compose up
