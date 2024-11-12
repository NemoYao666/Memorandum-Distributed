package mq

import (
	"strings"

	"github.com/streadway/amqp"

	"micro-todoList-k8s/config"
)

var RabbitMq *amqp.Connection

func InitRabbitMQ() {
	connString := strings.Join([]string{config.RabbitMQ, "://", config.RabbitMQUser, ":", config.RabbitMQPassWord, "@", config.RabbitMQHost, ":", config.RabbitMQPort, "/"}, "")
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	RabbitMq = conn
}
