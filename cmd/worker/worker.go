package main

import "github.com/samuel032khoury/gopherfeed/internal/mq"

type rabbitmqConfig = mq.Config

type mailConfig struct {
	fromEmail string
	host      string
	port      int
	username  string
	password  string
}
