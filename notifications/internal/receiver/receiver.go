package reciever

import (
	"context"
	"log"
	"route256/notifications/internal/handler"
	"route256/notifications/internal/infrastructure/kafka"
	"sync"

	"github.com/Shopify/sarama"
)

type HandleFunc func(ctx context.Context, message *sarama.ConsumerMessage)

type KafkaReceiver struct {
	consumer *kafka.ConsumerGroup
	handler  *handler.Handler
}

func NewReceiver(handler *handler.Handler) *KafkaReceiver {

	consumerGroup := kafka.NewConsumerGroup(handler.Notify)

	out := KafkaReceiver{
		consumer: &consumerGroup,
		handler:  handler,
	}

	return &out
}

func (r *KafkaReceiver) Subscribe(ctx context.Context, kafkaClient sarama.ConsumerGroup, topic string) error {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims

			if err := kafkaClient.Consume(ctx, []string{topic}, r.consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	wg.Wait()

	return nil
}
