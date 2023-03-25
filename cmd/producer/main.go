package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/moguchev/kafka/kafka"
	"github.com/moguchev/kafka/orders/sender"
)

var brokers = []string{
	"127.0.0.1:9091",
	"127.0.0.1:9092",
	"127.0.0.1:9093",
}

func main() {
	producer, err := kafka.NewSyncProducer(brokers)
	if err != nil {
		log.Fatalln(err)
	}

	asyncProducer, err := kafka.NewAsyncProducer(brokers)
	if err != nil {
		log.Fatalln(err)
	}

	onSuccess := func(id string) {
		fmt.Println("order success", id)
	}
	onFailed := func(id string) {
		fmt.Println("order failed", id)
	}

	sender := sender.NewOrderSender(
		producer,
		asyncProducer,
		"orders",
		onSuccess, onFailed,
	)
	s := &Server{
		OrderSender: sender,
	}

	http.HandleFunc("/order", s.CreateOrder)
	http.HandleFunc("/v2/order", s.CreateOrderV2)

	log.Println("Listen port: 8080")
	if err := http.ListenAndServe(":8080", nil); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err.Error())
	}
}

type OrderSender interface {
	SendOrderID(id int64) error
	SendOrderIDAsync(id int64)
}

// server
type Server struct {
	OrderSender
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusNotFound)
		return
	}

	// models
	type Order struct {
		ID int64 `json:"id"`
	}

	var req Order
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.OrderSender.SendOrderID(req.ID); err != nil {
		http.Error(w, "send failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) CreateOrderV2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusNotFound)
		return
	}

	// models
	type Order struct {
		ID int64 `json:"id"`
	}

	var req Order
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.OrderSender.SendOrderIDAsync(req.ID)

	w.WriteHeader(http.StatusCreated)
}
