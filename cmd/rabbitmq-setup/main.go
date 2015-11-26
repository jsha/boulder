// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

// This command does a one-time setup of the RabbitMQ exchange and the Activity
// Monitor queue, suitable for setting up a dev environment or Travis.

import (
	"fmt"

	"github.com/letsencrypt/boulder/Godeps/_workspace/src/github.com/cactus/go-statsd-client/statsd"

	"github.com/letsencrypt/boulder/cmd"
	blog "github.com/letsencrypt/boulder/log"
	"github.com/letsencrypt/boulder/rpc"
)

// Constants for AMQP
const (
	monitorQueueName = "Monitor"
	amqpExchange     = "boulder"
	amqpExchangeType = "topic"
	amqpInternal     = false
	amqpDurable      = false
	amqpDeleteUnused = false
	amqpExclusive    = false
	amqpNoWait       = false
)

func setup(c cmd.Config, _ statsd.Statter, _ *blog.AuditLogger) {
	ch, err := rpc.AmqpChannel(c.ActivityMonitor.AMQP)
	cmd.FailOnError(err, "Could not connect to AMQP")

	err = ch.ExchangeDeclare(
		amqpExchange,
		amqpExchangeType,
		amqpDurable,
		amqpDeleteUnused,
		amqpInternal,
		amqpNoWait,
		nil)
	cmd.FailOnError(err, "Declaring exchange")

	_, err = ch.QueueDeclare(
		monitorQueueName,
		amqpDurable,
		amqpDeleteUnused,
		amqpExclusive,
		amqpNoWait,
		nil)
	if err != nil {
		cmd.FailOnError(err, "Could not declare queue")
	}

	routingKey := "#" //wildcard

	err = ch.QueueBind(
		monitorQueueName,
		routingKey,
		amqpExchange,
		false,
		nil)
	if err != nil {
		txt := fmt.Sprintf("Could not bind to queue [%s]. NOTE: You may need to delete %s to re-trigger the bind attempt after fixing permissions, or manually bind the queue to %s.", monitorQueueName, monitorQueueName, routingKey)
		cmd.FailOnError(err, txt)
	}
}

func main() {
	app := cmd.NewAppShell("rabbitmq-setup", "Setup RabbitMQ")

	app.Action = setup
	app.Run()
}
