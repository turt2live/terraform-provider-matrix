#!/usr/bin/env sh

echo "starting synapse"
python /start.py &
SYNAPSE_PID=$!

echo "waiting for synapse to start"
RETRIES=60
while [ "$RETRIES" -gt 0 ] ; do
    sleep 1
    curl -sS --connect-timeout 60 http://localhost:8008/_matrix/client/versions && break
    let RETRIES=RETRIES-1
done

echo "creating admin account"
register_new_matrix_user -u admin -p test1234 -a -c /compiled/homeserver.yaml http://localhost:8008
access_token=$(curl -s -H 'Content-Type: application/json' --data '{"type":"m.login.password","user":"admin","password":"test1234"}' http://localhost:8008/_matrix/client/r0/login | jq .access_token | tr -d '"')
export MATRIX_ADMIN_ACCESS_TOKEN=$access_token
export MATRIX_CLIENT_SERVER_URL="http://localhost:8008"

echo "preparing project"
cd /project/src/github.com/turt2live/terraform-provider-matrix
export TF_ACC=true
#rm -rf vendor
#dep ensure -v -vendor-only

echo "running tests"
go test -v github.com/turt2live/terraform-provider-matrix/matrix
EXIT_CODE=$?

echo "killing synapse"
kill -9 $SYNAPSE_PID

echo "done (exit code $EXIT_CODE)"
exit $EXIT_CODE