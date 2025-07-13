package queue

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisEnvelope struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type RedisQueue struct {
	client *redis.Client
	key    string
	ctx    context.Context
	cancel context.CancelFunc
}

func NewRedisQueue(client *redis.Client, key string) *RedisQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &RedisQueue{
		client: client,
		key:    key,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (r *RedisQueue) Enqueue(jobType string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	env := redisEnvelope{
		Type: jobType,
		Data: data,
	}

	data, err = json.Marshal(env)
	if err != nil {
		return err
	}

	return r.client.LPush(r.ctx, r.key, data).Err()
}

func (r *RedisQueue) GetJobs() <-chan JobRunner {
	out := make(chan JobRunner)

	go func() {
		defer close(out)

		for {
			select {
			case <-r.ctx.Done():
				return
			default:
				res, err := r.client.BRPop(r.ctx, 5*time.Second, r.key).Result()
				if err != nil || len(res) < 2 {
					continue
				}

				var envelope redisEnvelope
				if err := json.Unmarshal([]byte(res[1]), &envelope); err != nil {
					log.Printf("invalid job envelope: %v", err)
					continue
				}

				job, ok := GetJob(envelope.Type)
				if !ok {
					log.Printf("unknown job type: %s", envelope.Type)
					continue
				}

				out <- JobRunner{Job: job, Data: envelope.Data}
			}
		}
	}()

	return out
}

func (r *RedisQueue) Close() {
	r.cancel()
}
