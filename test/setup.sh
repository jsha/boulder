#!/bin/bash
#
# Fetch dependencies of Boulderthat are necessary for development or testing,
# and configure database and RabbitMQ.
#

set -ev

go get \
  bitbucket.org/liamstask/goose/cmd/goose \
  golang.org/x/tools/cover \
  github.com/golang/lint/golint \
  github.com/tools/godep \
  github.com/mattn/goveralls \
  github.com/modocache/gover \
  github.com/jcjones/github-pr-status \
  github.com/jsha/listenbuddy &

# Set up rabbitmq exchange
go run cmd/rabbitmq-setup/main.go -server amqp://boulder-rabbitmq &

# Wait for all the background commands to finish.
wait
