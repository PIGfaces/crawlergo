package main

import (
	"testing"
	"time"

	"github.com/PIGfaces/crawlergo/pkg"
	gomock "github.com/golang/mock/gomock"

	"github.com/panjf2000/ants/v2"
	"github.com/stretchr/testify/assert"

	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"
)

func beforeTest(t *testing.T, mockProcess ProcessImp, poolNum, tabNum, tabTTLNum int) *pkg.CrawlerTask {
	pool, err := ants.NewPool(poolNum)
	assert.Nil(t, err)
	for i := 0; i <= poolNum; i++ {
		pool.Submit(func() {
			time.Sleep(time.Second * 5)
		})
	}
	task := pkg.CrawlerTask{
		Pool: pool,
		TabRunMonitor: pkg.TabRunMonitor{
			TabNum:    100,
			TabTTLNum: 50,
		},
		Config: &taskPkg.TaskConfig{
			MaxTabsCount: poolNum,
		},
	}
	ModifyPool(&task, mockProcess)
	return &task
}

func TestModifyPool_scale(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockProcess := NewMockProcessImp(mockCtl)
	mockProcess.EXPECT().CPUPercent().Return(0.1, nil)
	mockProcess.EXPECT().MemoryPercent().Return(float32(0.1), nil)
	// }()
	task := beforeTest(t, mockProcess, 10, 100, 50)
	assert.Equal(t, 17, task.Pool.Cap())
}

// 在最差的机器上
func TestModifyPool_reduce_half(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockProcess := NewMockProcessImp(mockCtl)
	mockProcess.EXPECT().CPUPercent().Return(0.8, nil)
	mockProcess.EXPECT().MemoryPercent().Return(float32(0.8), nil)

	task := beforeTest(t, mockProcess, 10, 100, 50)
	assert.Equal(t, 5, task.Pool.Cap())
}

func TestModifyPool_reduce_normal(t *testing.T) {
	mockCtl := gomock.NewController(t)
	mockProcess := NewMockProcessImp(mockCtl)
	mockProcess.EXPECT().CPUPercent().Return(0.8, nil)
	mockProcess.EXPECT().MemoryPercent().Return(float32(0.8), nil)

	task := beforeTest(t, mockProcess, 100, 100, 50)
	assert.Equal(t, 78, task.Pool.Cap())
}

func TestXxx(t *testing.T) {
	scaleJson := "{\"CpuWeight\": 0.3, \"MemWeight\": 0.8, \"TabTTLWeight\": 0.2}"
	ModifyScaleWeight(scaleJson)
	// destWeight := ModifyScaleWeight(scaleJson)
	assert.Equal(t, 0.3, sw.CpuWeight)
	assert.Equal(t, 0.8, sw.MemWeight)
	assert.Equal(t, 0.2, sw.TabTTLWeight)
}
