package sender

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/moguchev/kafka/pkg/proto/order"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type orderSender struct {
	producer      sarama.SyncProducer
	asyncProducer sarama.AsyncProducer
	topic         string
}

type Handler func(id string)

func NewOrderSender(producer1 sarama.SyncProducer, producer2 sarama.AsyncProducer, topic string, onSuccess, onFailed Handler) *orderSender {
	s := &orderSender{
		producer:      producer1,
		asyncProducer: producer2,
		topic:         topic,
	}

	// config.Producer.Return.Errors = true
	go func() {
		for e := range producer2.Errors() {
			bytes, _ := e.Msg.Key.Encode()

			onFailed(string(bytes))
			fmt.Println(e.Msg.Key, e.Error())
		}
	}()

	// config.Producer.Return.Successes = true
	go func() {
		for m := range producer2.Successes() {
			bytes, _ := m.Key.Encode()

			onSuccess(string(bytes))
			log.Printf("order id: %s, partition: %d, offset: %d", string(bytes), m.Partition, m.Offset)
		}
	}()

	return s
}

func (s *orderSender) SendOrderID(id int64) error {
	msg := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: -1,
		Value:     sarama.StringEncoder(fmt.Sprintf(`{"id":%d}`, id)),
		Key:       sarama.StringEncoder(fmt.Sprint(id)),
		Timestamp: time.Now(),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header-key"),
				Value: []byte("test-header-value"),
			},
		},
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("order id: %d, partition: %d, offset: %d", id, partition, offset)
	return nil
}

func (s *orderSender) SendOrderIDAsync(id int64) {
	// variant 1 (json)
	type Order struct {
		ID int64 `json:"id"`
	}
	bytes, err := json.Marshal(Order{
		ID: id,
	})
	if err != nil {
		return
	}

	// variant 2 (protobuf)
	orderpb := &order.Order{
		Id: id,
	}
	bytes, err = proto.Marshal(orderpb)
	if err != nil {
		return
	}
	// variant 3 (proto->json)
	bytes, err = protojson.Marshal(orderpb)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic:     s.topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(bytes),
		Key:       sarama.StringEncoder(fmt.Sprint(id)),
		Timestamp: time.Now(),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("test-header-key"),
				Value: []byte("test-header-value"),
			},
		},
	}

	s.asyncProducer.Input() <- msg
}
