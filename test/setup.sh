#!/bin/bash

go get \
  github.com/golang/tools/cmd/vet \
  github.com/golang/tools/cmd/cover \
  github.com/golang/lint/golint \
  github.com/mattn/goveralls \
  github.com/modocache/gover \
  github.com/jcjones/github-pr-status \
  github.com/jsha/listenbuddy &

(wget https://github.com/jsha/boulder-tools/raw/master/goose.gz &&
 mkdir -p $GOPATH/bin &&
 zcat goose.gz > $GOPATH/bin/goose &&
 chmod +x $GOPATH/bin/goose) &

# Set up rabbitmq exchange and activity monitor queue
go run cmd/rabbitmq-setup/main.go -server amqp://localhost &

./test/create_db.sh &

# Wait for all the background commands to finish.
wait
