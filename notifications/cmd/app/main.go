package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"route256/libs/logger"
	"route256/libs/srvwrapper"
	"route256/libs/tracer"
	"route256/libs/tx"
	"route256/notifications/internal/cash/cash"
	"route256/notifications/internal/config"
	configApp "route256/notifications/internal/config"
	"route256/notifications/internal/domain"
	"route256/notifications/internal/handlers/grpc"
	"route256/notifications/internal/handlers/http/getNotificationHistory"
	"route256/notifications/internal/infrastructure/kafka"
	"route256/notifications/internal/receiver"
	"route256/notifications/internal/repository/postgres"
	"route256/notifications/internal/telegram"
	"syscall"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	environment = flag.String("environment", "DEVELOPMENT", "environment: [DEVELOPMENT, PRODUCTION]")
)

const port = ":8089"

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

	BDPath := config.AppConfig.DSN()

	pool, err := pgxpool.Connect(ctx, BDPath)
	if err != nil {
		logger.Fatalf("connect to db: %s", err)
	}

	provider := tx.New(pool)
	repo := postgres.New(provider)

	rdb := redis.NewClient(&redis.Options{
		Addr:     configApp.AppConfig.Redis.Host,
		Password: configApp.AppConfig.Redis.Pass, // no password set
		DB:       configApp.AppConfig.Redis.DB,   // use default DB
	})

	cdb := cash.New(*rdb)

	m := domain.New(ctx, s, repo, cdb)

	consumerGroup, err := kafka.ConcumerGroupRun(configApp.AppConfig)

	if err != nil {
		logger.Fatalf("error in NewKafkaConsumerFromConfig: %w", err)
	}

	defer func() {
		if err = consumerGroup.Close(); err != nil {
			logger.Fatalf("Error in consumerGroup.Close(): %w", err)
		}
	}()

	hand := grpcHandler.NewHandler(m)

	res := reciever.NewReceiver(hand)

	getNotificationHistory := getNotificationHistory.Handler{Model: m}
	http.Handle("/getNotificationHistory", srvwrapper.New(getNotificationHistory.Handle))

	go func() {

		err = http.ListenAndServe(port, nil)
		if err != nil {
			logger.Fatalf("ERR: ", err)
		}

	}()

	err = res.Subscribe(ctx, consumerGroup, configApp.AppConfig.Kafka.TopicStatus)

	if err != nil {
		logger.Fatalf("Error in res.Subscribe: %w", err)
	}

}
