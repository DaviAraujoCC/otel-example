IMG ?= david13356/products-api
TARGET ?= development
PWD=$(shell pwd)
PORT=8080
INTERACTIVE:=$(shell [ -t 0 ] && echo i || echo d)
MAKE_FILE := $(lastword $(MAKEFILE_LIST))
PROJECT_NAME=products-api
MYSQL_USER ?= root
MYSQL_PASSWORD ?= root
MYSQL_HOST ?= localhost
MYSQL_PORT ?= 5432

MYSQL_DSN ?= $(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_HOST):$(MYSQL_PORT))


run: docker-build up-mysql
	@echo 'Running on http://localhost:$(PORT)/healthcheck'
	@docker run -t${INTERACTIVE} --rm \
		-v ${PWD}:/usr/app:delegated \
		-w /usr/app \
		--expose 8080 \
		-p $(PORT):8080 \
		--name ${PROJECT_NAME} \
		${IMG}
	@down-mysql
	
docker-build:
	@echo "Building Docker image..."
	@docker build --target ${TARGET} -t ${IMG} .

local-db:
	@docker run -d --rm \
		-e MYSQL_ROOT_PASSWORD=root \
		-p 3306:3306 \
		--name ${PROJECT_NAME}-mysql \
		mysql:5.7
	@echo "Waiting for MySQL to start..."
	@sleep 10
	@$(MAKE) -f $(MAKE_FILE) migrate-up

drop-db:
	@docker stop ${PROJECT_NAME}-mysql


migrate-setup:
	@if [ -z "$$(which migrate)" ]; then echo "Installing migrate command..."; go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate; fi


migrate-up: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path internal/db/migrations up $(N)


migrate-down: migrate-setup
	@ migrate -database 'mysql://$(MYSQL_DSN)?multiStatements=true' -path internal/db/migrations down $(N)

