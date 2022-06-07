package conf_test

import (
	"crawlergo/internal/conf"
	taskPkg "crawlergo/pkg/task"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func cmpTest(t *testing.T, cryptstr string, basicCmp *conf.RedisConf) {
	taskConf := taskPkg.TaskConfig{
		RedisConnectInfo: cryptstr,
	}
	testDest := conf.NewRedisConf(taskConf)
	assert.NotNil(t, testDest)
	assert.Equal(t, basicCmp, testDest, "not equal expacted")
}

func TestNewRedisConf(t *testing.T) {
	baseCryptInfo := "eyJjb25uZWN0aW9uIjogeyJob3N0IjogIjEyNy4wLjAuMSIsICJwb3J0IjogNjM3OSwgInBhc3N3b3JkIjogImxvY2FsdGVzdCJ9LCAidGFyZ2V0X2RiIjogMSwgInRhcmdldF9rZXkiOiAiYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoiLCAicmVzdWx0X2RiIjogMiwgImV4cGlyZWRfdGltZSI6IDEyMH0K"
	taskConf := taskPkg.TaskConfig{
		RedisConnectInfo: baseCryptInfo,
	}
	basicResult := &conf.RedisConf{
		Connection: conf.Connection{
			Host:     "127.0.0.1",
			Port:     6379,
			Password: "localtest",
		},
		TargetDB:        1,
		TargetKey:       "abcdefghijklmnopqrstuvwxyz",
		ResultDB:        2,
		ExpiredDuration: time.Minute * 2,
	}
	testDest := conf.NewRedisConf(taskConf)
	assert.NotNil(t, testDest)
	assert.Equal(t, basicResult, testDest, "not equal expacted")
}

func TestNewRedisConf_NotTTL(t *testing.T) {
	baseCryptInfo := "eyJjb25uZWN0aW9uIjogeyJob3N0IjogIjEyNy4wLjAuMSIsICJwb3J0IjogNjM3OSwgInBhc3N3b3JkIjogImxvY2FsdGVzdCJ9LCAidGFyZ2V0X2RiIjogMSwgInRhcmdldF9rZXkiOiAiYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXoiLCAicmVzdWx0X2RiIjogMn0K"
	basicResult := &conf.RedisConf{
		Connection: conf.Connection{
			Host:     "127.0.0.1",
			Port:     6379,
			Password: "localtest",
		},
		TargetDB:        1,
		TargetKey:       "abcdefghijklmnopqrstuvwxyz",
		ResultDB:        2,
		ExpiredDuration: time.Minute * 3,
	}
	cmpTest(t, baseCryptInfo, basicResult)
}
