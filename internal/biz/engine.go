package biz

import (
	"context"

	"github.com/PIGfaces/crawlergo/pkg/logger"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/PIGfaces/crawlergo/pkg/tools"
	"github.com/PIGfaces/crawlergo/pkg/tools/requests"

	"encoding/json"
	"fmt"

	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"
)

type Task struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type Tasks []*Task

func NewTask(id, url string) *Task {
	return &Task{
		ID:  id,
		Url: url,
	}
}

func NewTasks(taskInfo []byte) Tasks {
	var (
		taskMap = make(map[string]string)
		tasks   = make(Tasks, 0)
	)
	if err := json.Unmarshal(taskInfo, &taskMap); err != nil {
		logger.Logger.Fatal("cannot serialization taskInfo: ", string(taskInfo), " error: ", err.Error())
	}
	for key, val := range taskMap {
		tasks = append(tasks, NewTask(key, val))
	}
	return tasks
}

type EngineRepo interface {
	GetTaskValue(ctx context.Context) (Tasks, error)
	SetResult(ctx context.Context, key, value string) error
}

type EngineUsecase struct {
	repo        EngineRepo
	pushAddress string
}

func NewEngineUsecase(rep EngineRepo, taskConf taskPkg.TaskConfig) *EngineUsecase {
	return &EngineUsecase{
		repo:        rep,
		pushAddress: taskConf.PushAddress,
	}
}

func (eu *EngineUsecase) GetTasks() Tasks {
	tasks, err := eu.repo.GetTaskValue(context.Background())
	if err != nil {
		logger.Logger.Fatal("cannot get redis task", err.Error())
	}
	return tasks
}

func (eu *EngineUsecase) SetTaskResult(req *model.Request) {
	resp, err := requests.Request(req.Method, req.URL.String(), tools.ConvertHeaders(req.Headers), []byte(req.PostData), &requests.ReqOptions{
		Timeout: 1, AllowRedirect: false, Proxy: eu.pushAddress,
	})
	if err != nil {
		logger.Logger.Error("get html failed: ", err.Error())
		return
	}
	taskInfo := map[string]string{
		req.UniqueId(): req.URL.String(),
		"html":         resp.Text,
	}
	taskValue, _ := json.Marshal(taskInfo)

	taskKey := fmt.Sprintf("%s:%s", req.TaskID, req.UniqueId())
	if err = eu.repo.SetResult(context.Background(), taskKey, string(taskValue)); err != nil {
		logger.Logger.Error("set result failed, err: ", err.Error())
	}
}
