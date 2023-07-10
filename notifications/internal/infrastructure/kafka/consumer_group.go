package kafka

import (
	"context"
	"fmt"
	"log"
	"route256/notifications/internal/config"
	"time"

	"github.com/Shopify/sarama"
)

type ConsumerGroup struct {
	ready  chan bool
	Handle func(ctx context.Context, message *sarama.ConsumerMessage) error
}

func NewConsumerGroup(Handle func(ctx context.Context, message *sarama.ConsumerMessage) error) ConsumerGroup {
	return ConsumerGroup{
		ready:  make(chan bool),
		Handle: Handle}
}

func (consumer *ConsumerGroup) Ready() <-chan bool {
	return consumer.ready
}

// Setup Начинаем новую сессию, до ConsumeClaim
func (consumer *ConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	close(consumer.ready)

	return nil
}

// Cleanup завершает сессию, после того, как все ConsumeClaim завершатся
func (consumer *ConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читаем до тех пор пока сессия не завершилась
func (consumer *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {
		select {
		case message := <-claim.Messages():

			err := consumer.Handle(session.Context(), message)

			if err != nil {
				log.Printf("ERROR in ConsumeClaim.handle: %s", err)
			}

			// коммит сообщения "руками"
			session.MarkMessage(message, "")
		case <-session.Context().Done():
			return nil
		}
	}
}

func ConcumerGroupRun(AppConfig config.Config) (sarama.ConsumerGroup, error) {

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	configSarama := sarama.NewConfig()
	configSarama.Version = sarama.MaxVersion

	/*
		sarama.OffsetNewest - получаем только новые сообщений, те, которые уже были игнорируются
		sarama.OffsetOldest - читаем все с самого начала
	*/
	configSarama.Consumer.Offsets.Initial = sarama.OffsetOldest

	// Используется, если ваш offset "уехал" далеко и нужно пропустить невалидные сдвиги
	configSarama.Consumer.Group.ResetInvalidOffsets = true

	// Сердцебиение консьюмера
	configSarama.Consumer.Group.Heartbeat.Interval = 3 * time.Second

	// Таймаут сессии
	configSarama.Consumer.Group.Session.Timeout = 60 * time.Second

	// Таймаут ребалансировки
	configSarama.Consumer.Group.Rebalance.Timeout = 60 * time.Second

	const BalanceStrategy = "roundrobin"
	switch BalanceStrategy {
	case "sticky":
		configSarama.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategySticky}
	case "roundrobin":
		configSarama.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin}
	case "range":
		configSarama.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.BalanceStrategyRange}
	default:
		return nil, fmt.Errorf("Unrecognized consumer group partition assignor: %s", BalanceStrategy)
	}

	/**
	 * Setup a new Sarama consumer group
	 */

	client, err := sarama.NewConsumerGroup(AppConfig.Kafka.Brokers, AppConfig.Kafka.GroupName, configSarama)
	if err != nil {
		return nil, fmt.Errorf("Error creating consumer group client: %v", err)
	}

	return client, nil
}
