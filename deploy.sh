#!/bin/bash

# Service names
PRODUCT_SERVICE="product-service"
ORDER_SERVICE="order-service"
USER_SERVICE="user-service"
APP_STACK="app-stack"

# Docker image tag
TAG="latest"

# Kubernetes context (for kind)
KUBE_CONTEXT="kind-kind"

# Create kind cluster
create_kind_cluster() {
  echo "Creating a kind cluster..."
  kind create cluster --name kind --wait 60s
  kubectl cluster-info --context $KUBE_CONTEXT
}

# Build Docker images and load them into kind cluster
build_and_load_images() {
  echo "Building Docker images and loading into kind..."

  # Build Product Service Docker image
  cd $PRODUCT_SERVICE
  docker build -t $PRODUCT_SERVICE:$TAG .
  kind load docker-image $PRODUCT_SERVICE:$TAG --name kind
  cd ..

  # Build Order Service Docker image
  cd $ORDER_SERVICE
  docker build -t $ORDER_SERVICE:$TAG .
  kind load docker-image $ORDER_SERVICE:$TAG --name kind
  cd ..

  # Build User Service Docker image
  cd $USER_SERVICE
  docker build -t $USER_SERVICE:$TAG .
  kind load docker-image $USER_SERVICE:$TAG --name kind
  cd ..
}

# Deploy Product Service using Helm
deploy_product_service() {
  echo "Deploying Product Service..."
  helm upgrade -i $PRODUCT_SERVICE ./helm-charts/product-service \
    --namespace product-service \
    --set image.repository=$PRODUCT_SERVICE \
    --set image.tag=$TAG \
    --set database.enabled=true \
    --create-namespace
}

# Deploy Order Service using Helm
deploy_order_service() {
  echo "Deploying Order Service..."
  helm upgrade -i $ORDER_SERVICE ./helm-charts/order-service \
    --namespace order-service \
    --set image.repository=$ORDER_SERVICE \
    --set image.tag=$TAG \
    --set database.enabled=true \
    --create-namespace
}

# Deploy User Service using Helm
deploy_user_service() {
  echo "Deploying User Service..."
  helm upgrade -i $USER_SERVICE ./helm-charts/user-service \
    --namespace user-service \
    --set image.repository=$USER_SERVICE \
    --set image.tag=$TAG \
    --set database.enabled=true \
    --create-namespace
}

# Deploy All Service using Helm
deploy_all_services() {
  echo "Deploying All Services..."
  helm upgrade -i $APP_STACK ./helm-charts/all-services \
    --namespace all-services \
    --set user-service.image.repository=$USER_SERVICE \
    --set order-service.image.repository=$ORDER_SERVICE \
    --set product-service.image.repository=$PRODUCT_SERVICE \
    --set global.image.tag=$TAG \
    --create-namespace
}

# Update all chart deps in case charts were changed
update_chart_deps(){
  echo "Updating chart deps..."
  ./helm_charts/update.sh
}

# Main deployment process
main() {
  create_kind_cluster
  build_and_load_images
  update_chart_deps
  deploy_product_service
  deploy_order_service
  deploy_user_service
  deploy_all_services
}

# Execute main function
main
