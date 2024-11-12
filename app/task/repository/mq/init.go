package mq

import (
	"strings"

	"github.com/streadway/amqp"

	"micro-todoList-k8s/config"
)

var RabbitMq *amqp.Connection

func InitRabbitMQ() {
	connString := strings.Join([]string{config.C.RabbitMQ.RabbitMQ, "://", config.C.RabbitMQ.RabbitMQUser, ":", config.C.RabbitMQ.RabbitMQPassWord, "@", config.C.RabbitMQ.RabbitMQHost, ":", config.C.RabbitMQ.RabbitMQPort, "/"}, "")
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	RabbitMq = conn
}
