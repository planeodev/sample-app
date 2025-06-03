#!/bin/bash 

FOLDER=$(dirname "$0")

for d in user-service order-service product-service all-services ; do
    helm dependency update $FOLDER/$d
    helm dependency build $FOLDER/$d
done
