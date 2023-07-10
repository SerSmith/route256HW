package main

import (
	"context"
	"log"
	"os/signal"
	"route256/notifications/internal/config"
	configApp "route256/notifications/internal/config"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/handler"
	"route256/notifications/internal/infrastructure/kafka"
	"route256/notifications/internal/receiver"
	"route256/notifications/internal/telegram"
	"syscall"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := configApp.Init()

	if err != nil {
		log.Fatalln("error reading config: ", err)
	}

	log.Println("Config of APP: ", config.AppConfig)

	s, err := telegram.New(ctx, configApp.AppConfig.Telegram.Token, configApp.AppConfig.Telegram.Reciever_chat_id)

	if err != nil {
		log.Fatalln("error telegram.New: ", err)
	}

	m := domain.New(ctx, s)

	consumerGroup, err := kafka.ConcumerGroupRun(configApp.AppConfig)

	if err != nil {
		log.Fatalln("error in NewKafkaConsumerFromConfig: %w", err)
	}

	defer func() {
		if err = consumerGroup.Close(); err != nil {
			log.Fatalln("Error in consumerGroup.Close(): %w", err)
		}
	}()

	hand := handler.NewHandler(m)

	res := reciever.NewReceiver(hand)

	err = res.Subscribe(ctx, consumerGroup, configApp.AppConfig.Kafka.TopicStatus)

	if err != nil {
		log.Fatalln("Error in res.Subscribe: %w", err)
	}

}
