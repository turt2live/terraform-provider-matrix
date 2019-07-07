# terraform-provider-matrix

[![#terraform:t2bot.io](https://img.shields.io/badge/matrix-%23terraform:t2bot.io-brightgreen.svg)](https://matrix.to/#/#terraform:t2bot.io)
[![TravisCI badge](https://travis-ci.org/turt2live/terraform-provider-matrix.svg?branch=master)](https://travis-ci.org/turt2live/terraform-provider-matrix)
[![CircleCI](https://circleci.com/gh/turt2live/terraform-provider-matrix/tree/master.svg?style=svg)](https://circleci.com/gh/turt2live/terraform-provider-matrix/tree/master)
[![AppVeyor badge](https://ci.appveyor.com/api/projects/status/github/turt2live/terraform-provider-matrix?branch=master&svg=true)](https://ci.appveyor.com/project/turt2live/terraform-provider-matrix)

Terraform your matrix homeserver

## Building

Assuming Go 1.12 is already installed:
```bash
# Get it
git clone https://github.com/turt2live/terraform-provider-matrix
cd terraform-provider-matrix

# Install dependencies
go install -v ./...

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

## Usage

The matrix provider is a 3rd party plugin. See the documentation on [3rd party plugins](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins)
for installation instructions, and download the latest release from the [releases page](https://github.com/turt2live/terraform-provider-matrix/releases).

```hcl
provider "matrix" {
    # The client/server URL to access your matrix homeserver with.
    # Environment variable: MATRIX_CLIENT_SERVER_URL
    client_server_url = "https://matrix.org"
    
    # The default access token to use for things like content uploads.
    # Does not apply for provisioning users.
    # Environment variable: MATRIX_DEFAULT_ACCESS_TOKEN
    default_access_token = "MDAxSomeRandomString"
}
```

## Resources

The following resources are exposed from this provider.

### Media (Content)

Media (referred to as 'content' in the matrix specification) can be uploaded to the matrix content repository for later
use. Some uses include avatars for users, images in chat, etc. Media can also be existing before entering terraform and
referenced easily (skipping the upload process). Media cannot be deleted or updated.

Uploading media requires a `default_access_token` to be configured in the provider.

*Note*: Media cannot be deleted and is therefore abandoned when deleted in Terraform.

```hcl
# Existing media 
resource "matrix_content" "catpic" {
    # Your MXC URI must fit the following format/example: 
    #   Format:   mxc://origin/media_id
    #   Example:  mxc://matrix.org/SomeGeneratedId
    origin = "matrix.org"
    media_id = "SomeGeneratedId"
}

# New media (upload)
resource "matrix_content" "catpic" {
    file_path = "/path/to/cat_pic.png"
    file_name = "cat_pic.png"
    file_type = "image/png"
}
```

All media will have an `origin` and `media_id` as computed properties. To access the complete MXC URI, use the `id`.

### Users

Users can either be created using a username and password or by providing an access token. Users created with a username
and password will first be registered on the homeserver, and if the username appears to be in use then the provider will
try logging in.

*Note*: Users cannot be deleted and are therefore abandoned when deleted in Terraform.

```hcl
# Username/password user
resource "matrix_user" "foouser" {
    username = "foouser"
    password = "hunter2"
    
    # These properties are optional, and will update the user's profile
    # We're using a reference to the Media used in an earlier example
    display_name = "My Cool User"
    avatar_mxc = "${matrix_content.catpic.id}"
}

# Access token user
resource "matrix_user" "baruser" {
    access_token = "MDAxOtherCharactersHere"
    
    # These properties are optional, and will update the user's profile
    # We're using a reference to the Media used in an earlier example
    display_name = "My Cool User"
    avatar_mxc = "${matrix_content.catpic.id}"
}
```

All users have a `display_name`, `avatar_mxc`, and `access_token` as computed properties.

### Rooms

Rooms can be created by either specifying an explicit `room_id` or by specifying properties that help make up the room's
configuration for a new room. In both cases, a `member_access_token` is required because the provider needs an insight
into the room to perform state checks.

*Note*: Rooms cannot be completely deleted in matrix. When Terraform deletes a room, this provider will try to make the
room as inaccessible as possible. That generally means ensuring the `join_rules` are set to `private`, everyone is kicked,
aliases are removed, and the creator removes themselves from the room. For this reason, it is recommended that the member's
access token in the resource configuration be of at least power level 100 (Admin).

The examples here build off of previously mentioned resources, such as Users and Media.

```hcl
# Already existing room
resource "matrix_room" "fooroom" {
    room_id = "!somewhere:domain.com"
    member_access_token = "${matrix_user.foouser.access_token}"
}

# New room
resource "matrix_room" "barroom" {
    creator_user_id = "${matrix_user.foouser.id}"
    member_access_token = "${matrix_user.foouser.access_token}"
    
    # The rest is optional
    name = "My Room"
    avatar_mxc = "${matrix_content.catpic.id}"
    topic = "For testing only please"
    preset = "public_chat"
    guests_allowed = true
    invite_user_ids = ["${matrix_user.baruser.id}"]
    local_alias_localpart = "myroom"
}
```
