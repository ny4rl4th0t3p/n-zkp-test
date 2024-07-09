.PHONY: build up down

# Default goal
.DEFAULT_GOAL := build

build:
	@docker rmi n-zkp-test-verifier n-zkp-test-prover -f
	@docker build -t n-zkp-test-verifier -f Dockerfile.verifier .
	@docker build -t n-zkp-test-prover -f Dockerfile.prover .

up:
	$(MAKE) down
	@docker compose up -d
	@docker compose logs -f

down:
	@docker compose down