package main

type rabbitmqConfig struct {
	url       string
	queueName string
}

type mailConfig struct {
	fromEmail string
	host      string
	port      int
	username  string
	password  string
}
