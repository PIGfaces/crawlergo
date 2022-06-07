//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package pkg

import (
	"crawlergo/internal/biz"
	"crawlergo/internal/conf"
	"crawlergo/internal/data"

	taskPkg "crawlergo/pkg/task"
	"github.com/google/wire"
)

func wireRedisCacheRepo(taskPkg.TaskConfig) *biz.EngineUsecase {
	// func wireRedisCacheRepo(encryptInfo string) biz.EngineRepo {
	panic(wire.Build(conf.ProviderSet, data.ProviderSet, biz.ProviderSet))
	// panic(wire.Build(conf.ProviderSet, data.ProviderSet))
}
