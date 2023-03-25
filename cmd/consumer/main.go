package main

import (
	"context"
	"log"

	"github.com/moguchev/kafka/kafka"
	"github.com/moguchev/kafka/orders/reciever"
)

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func main() {
	c, err := kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatalln(err)
	}
	handlers := map[string]reciever.HandleFunc{
		"orders": func(id string) {
			log.Println(id)
		},
	}
	r := reciever.NewReciver(c, handlers)
	r.Subscribe("orders")

	<-context.TODO().Done()
}
