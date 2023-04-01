KIND_CLUSTER := otel-tracing
VERSION := 1.0

# ==============================================================================
# application build commands
tidy:
	go mod tidy
	go mod vendor

# ==============================================================================
# docker build
build-all: user-api sales-api product-api web-app

user-api:
	docker build \
		-f docker/dockerfile.user-api \
		-t user-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER \
		--build-arg VERSION=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

sales-api:
	docker build \
		-f docker/dockerfile.sales-api \
		-t sales-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER \
		--build-arg VERSION=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

product-api:
	docker build \
		-f docker/dockerfile.product-api \
		-t product-api-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER \
		--build-arg VERSION=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

web-app:
	docker build \
		-f docker/dockerfile.web-app \
		-t web-app-amd64:${VERSION} \
		--build-arg BUILD_ENV=DOCKER \
		--build-arg VERSION=${VERSION} \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# kind cluster
kind-up:
	kind create cluster --name ${KIND_CLUSTER} --config k8s/kind/config.yaml
	kubectl create namespace webapp-ns
	kubectl create namespace userapi-ns
	kubectl create namespace salesapi-ns
	kubectl create namespace productapi-ns
	kubectl config set-context --current --namespace webapp-ns

kind-deploy: build-all kind-load kind-apply

kind-update: kind-load kind-restart

kind-load: kind-load-user-api kind-load-sales-api kind-load-product-api kind-load-web-app

kind-load-user-api:
	cd k8s/userapi/overlays/kind && kustomize edit set image user-api-image=user-api-amd64:${VERSION}
	kind load docker-image user-api-amd64:${VERSION} --name ${KIND_CLUSTER}

kind-load-sales-api:
	cd k8s/salesapi/overlays/kind && kustomize edit set image sales-api-image=sales-api-amd64:${VERSION}
	kind load docker-image sales-api-amd64:${VERSION} --name ${KIND_CLUSTER}

kind-load-product-api:
	cd k8s/productapi/overlays/kind && kustomize edit set image product-api-image=product-api-amd64:${VERSION}
	kind load docker-image product-api-amd64:${VERSION} --name ${KIND_CLUSTER}

kind-load-web-app:
	cd k8s/webapp/overlays/kind && kustomize edit set image web-app-image=web-app-amd64:${VERSION}
	kind load docker-image web-app-amd64:${VERSION} --name ${KIND_CLUSTER}

kind-deploy-zipkin:
	kustomize build k8s/zipkin/ | kubectl apply -f -
	kubectl wait --namespace tracer-ns --timeout=120s --for=condition=Available deployment/zipkin

kind-apply:
	kustomize build k8s/productapi/overlays/kind/ | kubectl apply -f -
	kubectl wait --namespace productapi-ns --timeout=120s --for=condition=Available deployment/product-api

	kustomize build k8s/salesapi/overlays/kind/ | kubectl apply -f -	
	kubectl wait --namespace salesapi-ns --timeout=120s --for=condition=Available deployment/sales-api

	kustomize build k8s/userapi/overlays/kind/ | kubectl apply -f -
	kubectl wait --namespace userapi-ns --timeout=120s --for=condition=Available deployment/user-api

	kustomize build k8s/webapp/overlays/kind/ | kubectl apply -f -
	kubectl wait --namespace webapp-ns --timeout=120s --for=condition=Available deployment/web-app

kind-delete-deploy:
	kustomize build k8s/productapi/overlays/kind/ | kubectl delete -f -
	kustomize build k8s/salesapi/overlays/kind/ | kubectl delete -f -
	kustomize build k8s/userapi/overlays/kind/ | kubectl delete -f -
	kustomize build k8s/webapp/overlays/kind/ | kubectl delete -f -

kind-restart:
	kubectl rollout restart deployment user-api --namespace userapi-ns
	kubectl rollout restart deployment product-api --namespace productapi-ns
	kubectl rollout restart deployment sales-api --namespace salesapi-ns
	kubectl rollout restart deployment web-app --namespace webapp-ns

kind-down:
	kind delete cluster --name ${KIND_CLUSTER}