package main

import (
	"context"
	"flag"
	"os/signal"
	"route256/libs/logger"
	"route256/libs/tracer"
	"route256/notifications/internal/config"
	configApp "route256/notifications/internal/config"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/handler"
	"route256/notifications/internal/infrastructure/kafka"
	"route256/notifications/internal/receiver"
	"route256/notifications/internal/telegram"
	"syscall"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.SetLoggerByEnvironment(*environment)

	if err := tracer.InitGlobal(domain.ServiceName); err != nil {
		logger.Fatalf("error init tracer: ", err)
	}

	err := configApp.Init()

	if err != nil {
		logger.Fatalf("error reading config: ", err)
	}

	logger.Info("Config of APP: ", config.AppConfig)

	s, err := telegram.New(ctx, configApp.AppConfig.Telegram.Token, configApp.AppConfig.Telegram.Reciever_chat_id)

	if err != nil {
		logger.Fatalf("error telegram.New: ", err)
	}

	m := domain.New(ctx, s)

	consumerGroup, err := kafka.ConcumerGroupRun(configApp.AppConfig)

	if err != nil {
		logger.Fatalf("error in NewKafkaConsumerFromConfig: %w", err)
	}

	defer func() {
		if err = consumerGroup.Close(); err != nil {
			logger.Fatalf("Error in consumerGroup.Close(): %w", err)
		}
	}()

	hand := handler.NewHandler(m)

	res := reciever.NewReceiver(hand)

	err = res.Subscribe(ctx, consumerGroup, configApp.AppConfig.Kafka.TopicStatus)

	if err != nil {
		logger.Fatalf("Error in res.Subscribe: %w", err)
	}

}
