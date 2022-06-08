//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package pkg

import (
	"github.com/PIGfaces/crawlergo/internal/biz"
	"github.com/PIGfaces/crawlergo/internal/conf"
	"github.com/PIGfaces/crawlergo/internal/data"

	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"
	"github.com/google/wire"
)

func wireRedisCacheRepo(taskPkg.TaskConfig) *biz.EngineUsecase {
	// func wireRedisCacheRepo(encryptInfo string) biz.EngineRepo {
	panic(wire.Build(conf.ProviderSet, data.ProviderSet, biz.ProviderSet))
	// panic(wire.Build(conf.ProviderSet, data.ProviderSet))
}
