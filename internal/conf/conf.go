package conf

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/PIGfaces/crawlergo/pkg/logger"

	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"

	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewRedisConf,
)

type RedisConf struct {
	Connection      Connection    `json:"connection"`
	TargetDB        int           `json:"target_db"`
	TargetKey       string        `json:"target_key"`
	ResultDB        int           `json:"result_db"`
	ExpiredDuration time.Duration `json:"expired_time"`
}

type Connection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

func NewRedisConf(taskConf taskPkg.TaskConfig) *RedisConf {
	connInfo, err := base64.StdEncoding.DecodeString(taskConf.RedisConnectInfo)
	if err != nil {
		logger.Logger.Fatal("base64 decode encrypt string failed, error info: ", err.Error(), " encryptInfo: ", taskConf.RedisConnectInfo)
	}

	var rdbConf RedisConf
	if err = json.Unmarshal(connInfo, &rdbConf); err != nil {
		logger.Logger.Fatal("unserialization failed: ", err.Error(), " encryptInfo: ", taskConf.RedisConnectInfo)
	}
	if rdbConf.ExpiredDuration == 0 {
		// 默认为3分钟
		rdbConf.ExpiredDuration = time.Minute * 3
	} else {
		rdbConf.ExpiredDuration *= time.Second
	}
	return &rdbConf
}
