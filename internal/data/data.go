package data

import (
	"crawlergo/internal/conf"
	"fmt"

	"github.com/go-redis/redis/extra/redisotel"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewEngineRepo,
)

type (
	Data struct {
		readRdb  *redis.Client
		writeRdb *redis.Client
	}

	DataOptFunc func(*Data)
)

func newRDB(host, pwd string, db int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: pwd,
		DB:       db,
		// DB:           int(conf.Redis.Db),
		// DialTimeout: conf.Redis.DialTimeout.AsDuration(),
		// WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		// ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	rdb.AddHook(redisotel.TracingHook{})
	return rdb
}

func NewData(conf *conf.RedisConf) *Data {
	data := new(Data)
	// 设置 data 的属性,
	host := fmt.Sprintf("%s:%d", conf.Connection.Host, conf.Connection.Port)
	for _, fn := range []DataOptFunc{
		WithReadDB(host, conf.Connection.Password, conf.TargetDB),
		WithWriteRDB(host, conf.Connection.Password, conf.ResultDB),
		// WithTargetKey(conf.TargetKey),
	} {
		fn(data)
	}
	return data
}

func WithWriteRDB(host, pwd string, db int) DataOptFunc {
	return func(d *Data) {
		d.writeRdb = newRDB(host, pwd, db)
	}
}

// TargetKey 从 redis 中读取任务的 key
// func WithTargetKey(readKey string) DataOptFunc {
// 	return func(d *Data) {
// 		d.targetKey = readKey
// 	}
// }

func WithReadDB(host, pwd string, db int) DataOptFunc {
	return func(d *Data) {
		d.readRdb = newRDB(host, pwd, db)
	}
}
