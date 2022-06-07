package data

import (
	"context"
	"crawlergo/internal/biz"
	"crawlergo/internal/conf"
	"time"
)

type EngineRepo struct {
	data      *Data
	targetKey string
	cacheTTl  time.Duration
}

func NewEngineRepo(data *Data, c *conf.RedisConf) biz.EngineRepo {
	return &EngineRepo{
		data:      data,
		targetKey: c.TargetKey,
		cacheTTl:  c.ExpiredDuration,
	}
}

func (er *EngineRepo) GetTaskValue(ctx context.Context) (biz.Tasks, error) {
	value, err := er.data.readRdb.Get(ctx, er.targetKey).Bytes()
	if err != nil {
		return nil, err
	}
	return biz.NewTasks(value), nil
}

func (er *EngineRepo) SetResult(ctx context.Context, key, value string) error {
	return er.data.writeRdb.Set(ctx, key, value, er.cacheTTl).Err()
}

func (er *EngineRepo) GetResult(ctx context.Context, key string) ([]byte, error) {
	return er.data.writeRdb.Get(ctx, key).Bytes()
}

func (er *EngineRepo) SetRead(ctx context.Context, key, value string) error {
	return er.data.readRdb.Set(ctx, key, value, er.cacheTTl).Err()
}
