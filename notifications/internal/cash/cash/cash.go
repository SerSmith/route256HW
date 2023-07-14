package cash

import (
	"context"
	"encoding/json"
	"fmt"
	"route256/libs/logger"
	"route256/libs/tracer"
	"route256/notifications/internal/cash/converter"
	"route256/notifications/internal/cash/schema"
	"route256/notifications/internal/domain"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/redis/go-redis/v9"
)

type cashDB struct {
	redis *redis.Client
}

func New(redis redis.Client) *cashDB {
	return &cashDB{&redis}
}

func (c *cashDB) Get(ctx context.Context, req domain.NotificationHistoryRequest) ([]domain.NotificationMem, bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cashDB/Get")
	defer span.Finish()

	reqSchema := converter.NotificationHistoryRequestD2S(req)

	key, err := json.Marshal(reqSchema)
	if err != nil {
		return nil, false, tracer.MarkSpanWithError(ctx, err)
	}

	quantStr, err := c.redis.Get(ctx, c.constructKeyQuant(ctx, string(key))).Result()

	if err == redis.Nil {
		return nil, false, nil
	} else if err != nil {
		return nil, false, tracer.MarkSpanWithError(ctx, err)
	}

	quant, err := strconv.Atoi(quantStr)

	if err != nil {
		return nil, false, tracer.MarkSpanWithError(ctx, err)
	}

	outSchema := make([]schema.StatusChangeMessage, 0, quant)

	for i := 0; i < quant; i++ {
		keyElement := c.constructKeyElement(ctx, string(key), i)

		val, err := c.redis.Get(ctx, keyElement).Result()

		if err != nil {
			return nil, false, tracer.MarkSpanWithError(ctx, err)
		}

		scmSchema := schema.StatusChangeMessage{}
		err = json.Unmarshal([]byte(val), &scmSchema)

		if err != nil {
			return nil, false, tracer.MarkSpanWithError(ctx, err)
		}

		outSchema = append(outSchema, scmSchema)

	}

	out := converter.NotificationMemS2D(outSchema)
	return out, true, err

}

func (c *cashDB) Set(ctx context.Context, req domain.NotificationHistoryRequest, value []domain.NotificationMem) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "cashDB/Set")
	defer span.Finish()

	reqSchema := converter.NotificationHistoryRequestD2S(req)

	key, err := json.Marshal(reqSchema)
	if err != nil {
		return tracer.MarkSpanWithError(ctx, err)
	}

	valueSchema := converter.NotificationMemD2S(value)

	keyQuant := c.constructKeyQuant(ctx, string(key))
	c.redis.Set(ctx, keyQuant, strconv.Itoa(len(valueSchema)), redis.KeepTTL)

	for i, vs := range valueSchema {
		keyElement := c.constructKeyElement(ctx, string(key), i)

		vsMarshal, err := json.Marshal(vs)
		if err != nil {
			return tracer.MarkSpanWithError(ctx, err)
		}

		c.redis.Set(ctx, keyElement, string(vsMarshal), redis.KeepTTL)
	}

	return nil
}

func (c *cashDB) constructKeyQuant(ctx context.Context, key string) string {
	return fmt.Sprintf("%s_quant", string(key))
}

func (c *cashDB) constructKeyElement(ctx context.Context, key string, num int) string {
	return fmt.Sprintf("%s_%d", key, num)
}
