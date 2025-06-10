## Sample App

This example app represents three distinct services and associated databases. To install a specific service, such as the order service, use that helm chart. Be sure to pass values to enable the DB as needed.

To install all of the services in one go, as a combined umbrella chart, use the all-services chart which has its own shared DB.
 
Be sure to build and tag the images for each service and make them available in K8s before deploying

### Building images

Before you can run these services you'll need to build the images first

```
$ cd order-service && docker build -t order-service:latest . && cd ..
$ cd user-service && docker build -t user-service:latest . && cd ..
$ cd product-service && docker build -t product-service:latest . && cd ..
```

Keep in mind that this will provide those images locally, and your kubernetes cluster may also need access to them. This may require pushing them to a remote registry, or using another tool.

### Running a single service

To run a specific service, such as the order service, use that specific helm chat within the `helm-charts` folder. To run it with a working database we'll need to enable the database via the values file paramter, and also ensure the database host matches our generated service.

```
$ cd helm-charts/order-service
# first build dependencies
$ helm dependency build
# next install the helm chart, enabling the database
$ helm install myorder-service . -f values.yaml --set database.enabled=true --set global.postgresql.host=myorder-service-postgresql.
# forward the service port locally
$ kubectl port-forward svc/myorder-service-order-service 8088:8081
# access the running service
$ curl localhost:8088
# also hit the endpoint that accesses the DB
$ curl localhost:8088/orders
```

If everything is working you should see a DB schema error.

### Running all services

To run all the services with a shared DB, use the `all-services` helm chat within the `helm-charts` folder. To run it with a working with the included database, we'll need to ensure the database host matches our generated service.

```
$ cd helm-charts/all-services
# first build dependencies
$ helm dependency build
# next install the helm chart, ensuring the DB host matches
$ helm install myappstack . -f values.yaml --set global.postgresql.host=myappstack-postgresql
# forward the service port locally
$ kubectl port-forward svc/myappstack-order-service 8088:8081
# access the running service
$ curl localhost:8088
# also hit the endpoint that accesses the DB
$ curl localhost:8088/orders
```

If everything is working you should see a DB schema error.