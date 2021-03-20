SHELL = /bin/sh
DC_BIN = $(shell command -v docker-compose 2> /dev/null)
DC_RUN_ARGS = --rm --user "$(shell id -u):$(shell id -g)"

build: ## Build
	$(DC_BIN) build

up: ## Create and start containers
	$(DC_BIN) up --detach postgres

down: ## Create and start containers
	$(DC_BIN) down -t 5

shell: ## Start shell into app container
	$(DC_BIN) run $(DC_RUN_ARGS) app sh

pull: ## Pulling newer versions of used docker images
	$(DC_BIN) pull