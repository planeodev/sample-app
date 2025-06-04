## Sample App

This example app represents three distinct services and associated databases. To install a specific service, such as the order service, use that helm chart. Be sure to pass values to enable the DB as needed.

To install all of the services in one go, as a combined umbrella chart, use the all-services chart which has its own shared DB.
 
Be sure to build and tag the images for each service and make them available in K8s before deploying
