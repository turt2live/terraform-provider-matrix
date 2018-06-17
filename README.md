# terraform-provider-matrix

[![#terraform:t2bot.io](https://img.shields.io/badge/matrix-%23terraform:t2bot.io-brightgreen.svg)](https://matrix.to/#/#terraform:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/terraform-provider-matrix.svg?branch=master)](https://travis-ci.org/turt2live/terraform-provider-matrix)

Terraform your matrix homeserver

# Building

Assuming Go 1.9 and `dep` are already installed:
```bash
# Get it
git clone https://github.com/turt2live/terraform-provider-matrix
cd terraform-provider-matrix

# Grab the dependencies
dep ensure

# Build it
go build -o terraform-provider-matrix
```