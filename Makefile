NAME=otel-example
IMG?=docker-repository/${NAME}

build-app:
	@echo "Building..."
	@CGO_ENABLED=0 go build -o bin/$(NAME) -v

docker-build:
	@echo "Building docker image..."
	@docker build -t $(IMG) app/.

docker-push:
	@echo "Pushing docker image..."
	@docker push $(IMG)


k8s-deploy:
	@echo "Deploying to k8s..."
	@kubectl apply -k k8s/.

k8s-destroy:
	@echo "Destroying k8s deployment..."
	@kubectl delete -k k8s/.

up: docker-build docker-push k8s-deploy
	@echo "Starting..."
	
down: k8s-destroy
	@echo "Stopping..."