# terraform-provider-matrix

[![#terraform:t2bot.io](https://img.shields.io/badge/matrix-%23terraform:t2bot.io-brightgreen.svg)](https://matrix.to/#/#terraform:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/terraform-provider-matrix.svg?branch=master)](https://travis-ci.org/turt2live/terraform-provider-matrix)

Terraform your matrix homeserver

## Building

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

## Running the tests

The tests run within a Docker container. This is to ensure that the test homeserver gets set up correctly and doesn't 
leave lingering data on another homeserver.

The tests can be run with `./run-tests.sh` or by running the following commands:
```
docker build -t terraform-provider-matrix-tests .
docker run --name terraform-provider-matrix-tests terraform-provider-matrix-tests
docker rm terraform-provider-matrix-tests
```

The first execution may take a while to set up, however future executions should be
fairly quick.

## Resources

The following resources are exposed from this provider.

### Users

Users can either be created using a username and password or by providing an access token. Users created with a username
and password will first be registered on the homeserver, and if the username appears to be in use then the provider will
try logging in.

```hcl
# Username/password user
resource "matrix_user" "foouser" {
    username = "foouser"
    password = "hunter2"
}

# Access token user
resource "matrix_user" "baruser" {
    access_token = "MDAxOtherCharactersHere"
}
```

All users have a `display_name`, `avatar_mxc`, and `access_token` as computed properties.