#!/bin/bash

docker build -t terraform-provider-matrix-tests .
docker run --rm --name terraform-provider-matrix-tests -e TF_LOG=$TF_LOG terraform-provider-matrix-tests
