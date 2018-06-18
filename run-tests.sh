#!/bin/bash

docker build -t terraform-provider-matrix-tests .
docker run --name terraform-provider-matrix-tests terraform-provider-matrix-tests
docker rm terraform-provider-matrix-tests
