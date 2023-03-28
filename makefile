KIND_CLUSTER := otel-dist-tracing-cluster
VERSION := 1.0

# ==============================================================================
# application build commands
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# docker build
all: user-api sales-api product-api

user-api:
	docker build \
		-f docker/dockerfile.user-api \
		-t user-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER_${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

sales-api:
	docker build \
		-f docker/dockerfile.sales-api \
		-t sales-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER_${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

product-api:
	docker build \
		-f docker/dockerfile.product-api \
		-t product-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER_${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# kind cluster
kind-up:
	kind create cluster --name ${KIND_CLUSTER} --config k8s/kind/config.yaml
	kubectl create namespace usersalesapi-ns
	kubectl create namespace userapi-ns
	kubectl create namespace salesapi-ns
	kubectl create namespace productapi-ns
	kubectl config set-context --current --namespace usersalesapi-ns

kind-down:
	kind delete cluster --name ${KIND_CLUSTER}