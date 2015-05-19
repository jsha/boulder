#!/bin/bash
cd $(dirname $0)/..

kill_boulder() {
  fuser -skn tcp 4300
  fuser -skn tcp 9300
}
# Kill any leftover boulder or cfssl processes from previous runs.
kill_boulder

go run ./cmd/boulder/main.go --config test/boulder-test-config.json &>/dev/null &
go run Godeps/_workspace/src/github.com/cloudflare/cfssl/cmd/cfssl/cfssl.go \
  -loglevel 0 \
  serve \
  -port 9300 \
  -ca test/test-ca.pem \
  -ca-key test/test-ca.key \
  -config test/cfssl-config.json &>/dev/null &

cd test/js
npm install

# Wait for Boulder to come up
until nc localhost 4300 < /dev/null ; do sleep 1 ; done

js test.js --email foo@bar.com --agree true --domain foo.com --new-reg http://localhost:4300/acme/new-reg
STATUS=$?

kill_boulder
exit $STATUS
