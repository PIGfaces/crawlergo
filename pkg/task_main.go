package pkg

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/PIGfaces/crawlergo/internal/biz"
	"github.com/PIGfaces/crawlergo/pkg/config"
	engine2 "github.com/PIGfaces/crawlergo/pkg/engine"
	filter2 "github.com/PIGfaces/crawlergo/pkg/filter"
	"github.com/PIGfaces/crawlergo/pkg/logger"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/PIGfaces/crawlergo/pkg/resultsave"
	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"

	"github.com/panjf2000/ants/v2"
)

type CrawlerTask struct {
	Browser       *engine2.Browser //
	RootDomain    string           // 当前爬取根域名 用于子域名收集
	Targets       []*model.Request // 输入目标
	Result        *Result          // 最终结果
	ResultSave    resultsave.ResultSave
	Config        *taskPkg.TaskConfig // 配置信息
	smartFilter   filter2.SmartFilter // 过滤对象
	Pool          *ants.Pool          // 协程池
	taskWG        sync.WaitGroup      // 等待协程池所有任务结束
	crawledCount  int                 // 爬取过的数量
	taskCountLock sync.Mutex          // 已爬取的任务总数锁
	redisUsecase  *biz.EngineUsecase  // 跟 redis 交互的接口
}

type Result struct {
	AllReqSimpFilter *filter2.SimpleFilter
	ReqSave          resultsave.ResultSave
	AllReqSave       resultsave.ResultSave
	allDomainSave    resultsave.ResultSave
	subDomainSave    resultsave.ResultSave
	resultLock       sync.Mutex // 合并结果时加锁
}

type TaskConfig struct {
	MaxCrawlCount           int    // 最大爬取的数量
	FilterMode              string // simple、smart、strict
	ExtraHeaders            map[string]interface{}
	ExtraHeadersString      string
	AllDomainReturn         bool // 全部域名收集
	SubDomainReturn         bool // 子域名收集
	IncognitoContext        bool // 开启隐身模式
	NoHeadless              bool // headless模式
	DomContentLoadedTimeout time.Duration
	TabRunTimeout           time.Duration     // 单个标签页超时
	PathByFuzz              bool              // 通过字典进行Path Fuzz
	FuzzDictPath            string            //Fuzz目录字典
	PathFromRobots          bool              // 解析Robots文件找出路径
	MaxTabsCount            int               // 允许开启的最大标签页数量 即同时爬取的数量
	ChromiumPath            string            // Chromium的程序路径  `/home/zhusiyu1/chrome-linux/chrome`
	EventTriggerMode        string            // 事件触发的调用方式： 异步 或 顺序
	EventTriggerInterval    time.Duration     // 事件触发的间隔
	BeforeExitDelay         time.Duration     // 退出前的等待时间，等待DOM渲染，等待XHR发出捕获
	EncodeURLWithCharset    bool              // 使用检测到的字符集自动编码URL
	IgnoreKeywords          []string          // 忽略的关键字，匹配上之后将不再扫描且不发送请求
	Proxy                   string            // 请求代理
	CustomFormValues        map[string]string // 自定义表单填充参数
	CustomFormKeywordValues map[string]string // 自定义表单关键词填充内容
}

type tabTask struct {
	crawlerTask *CrawlerTask
	browser     *engine2.Browser
	req         *model.Request
}

/**
新建爬虫任务
*/
func NewCrawlerTask(urls []string, taskConf taskPkg.TaskConfig, postData string) (*CrawlerTask, error) {

	var (
		CSPEngineuc  *biz.EngineUsecase
		redisTargets []*model.Request
		taskResult   = &Result{
			AllReqSimpFilter: &filter2.SimpleFilter{},
		}
		targets = []*model.Request{}
	)
	// redis 中的任务
	if taskConf.RedisConnectInfo != "" {
		redisTargets, CSPEngineuc = GetRedisTarget(taskConf, postData)
		targets = append(targets, redisTargets...)
	}
	// 命令行的参数任务
	targets = append(targets, MakeTargets(urls, taskConf, postData)...)

	if taskConf.OutputJsonPath != "" {
		taskResult.AllReqSave = resultsave.NewFileSave(fmt.Sprintf("%s/%s", taskConf.OutputJsonPath, config.ALL_REQUEST_FILE))
		taskResult.ReqSave = resultsave.NewFileSave(fmt.Sprintf("%s/%s", taskConf.OutputJsonPath, config.REQUEST_FILE))
		taskResult.allDomainSave = resultsave.NewAllDomainSave(fmt.Sprintf("%s/%s", taskConf.OutputJsonPath, config.ALL_DOMAIN_FILE))
		taskResult.subDomainSave = resultsave.NewDomainSave(fmt.Sprintf("%s/%s", taskConf.OutputJsonPath, config.SUB_DOMAIN_FILE), targets[0].URL.RootDomain())
	}

	crawlerTask := CrawlerTask{
		Result: taskResult,
		Config: &taskConf,
		smartFilter: filter2.SmartFilter{
			SimpleFilter: filter2.SimpleFilter{
				HostLimit: targets[0].URL.Host,
			},
		},
		redisUsecase: CSPEngineuc,
	}

	if len(targets) == 1 {
		_newReq := *targets[0]
		newReq := &_newReq
		_newURL := *_newReq.URL
		newReq.URL = &_newURL
		if targets[0].URL.Scheme == "http" {
			newReq.URL.Scheme = "https"
		} else {
			newReq.URL.Scheme = "http"
		}
		targets = append(targets, newReq)
	}
	crawlerTask.Targets = targets[:]

	for _, req := range targets {
		req.Source = config.FromTarget
	}

	if len(targets) == 0 {
		logger.Logger.Fatal("no validate target.")
	}
	logger.Logger.Info(fmt.Sprintf("Init crawler task, host: %s, max tab count: %d, max crawl count: %d.",
		targets[0].URL.Host, taskConf.MaxTabsCount, taskConf.MaxCrawlCount))
	logger.Logger.Info("filter mode: ", taskConf.FilterMode)

	// 业务代码与数据代码分离, 初始化一些默认配置
	for _, fn := range []taskPkg.TaskConfigOptFunc{
		taskPkg.WithTabRunTimeout(config.TabRunTimeout),
		taskPkg.WithMaxTabsCount(config.MaxTabsCount),
		taskPkg.WithMaxCrawlCount(config.MaxCrawlCount),
		taskPkg.WithDomContentLoadedTimeout(config.DomContentLoadedTimeout),
		taskPkg.WithEventTriggerInterval(config.EventTriggerInterval),
		taskPkg.WithBeforeExitDelay(config.BeforeExitDelay),
		taskPkg.WithEventTriggerMode(config.DefaultEventTriggerMode),
		taskPkg.WithIgnoreKeywords(config.DefaultIgnoreKeywords),
	} {
		fn(&taskConf)
	}

	if taskConf.MaxCrawlCount < len(crawlerTask.Targets) {
		// 如果最大爬取数量都少于任务数量就会不完整了
		taskConf.MaxCrawlCount = len(crawlerTask.Targets) * 100
	}

	if taskConf.ExtraHeadersString != "" {
		err := json.Unmarshal([]byte(taskConf.ExtraHeadersString), &taskConf.ExtraHeaders)
		if err != nil {
			logger.Logger.Error("custom headers can't be Unmarshal.")
			return nil, err
		}
	}

	crawlerTask.Browser = engine2.InitBrowser(taskConf.ChromiumPath, taskConf.IncognitoContext, taskConf.ExtraHeaders, taskConf.Proxy, taskConf.NoHeadless)
	crawlerTask.RootDomain = targets[0].URL.RootDomain()

	crawlerTask.smartFilter.Init()

	// 创建协程池
	p, _ := ants.NewPool(taskConf.MaxTabsCount)
	crawlerTask.Pool = p

	return &crawlerTask, nil
}

/**
根据请求列表生成tabTask协程任务列表
*/
func (t *CrawlerTask) generateTabTask(req *model.Request) *tabTask {
	task := tabTask{
		crawlerTask: t,
		browser:     t.Browser,
		req:         req,
	}
	return &task
}

/**
开始当前任务
*/
func (t *CrawlerTask) Run() {
	defer t.Pool.Release()  // 释放协程池
	defer t.Browser.Close() // 关闭浏览器

	if t.Config.PathFromRobots {
		reqsFromRobots := GetPathsFromRobots(*t.Targets[0])
		logger.Logger.Info("get paths from robots.txt: ", len(reqsFromRobots))
		t.Targets = append(t.Targets, reqsFromRobots...)
	}

	if t.Config.FuzzDictPath != "" {
		if t.Config.PathByFuzz {
			logger.Logger.Warn("`--fuzz-path` is ignored, using `--fuzz-path-dict` instead")
		}
		reqsByFuzz := GetPathsByFuzzDict(*t.Targets[0], t.Config.FuzzDictPath)
		t.Targets = append(t.Targets, reqsByFuzz...)
	} else if t.Config.PathByFuzz {
		reqsByFuzz := GetPathsByFuzz(*t.Targets[0])
		logger.Logger.Info("get paths by fuzzing: ", len(reqsByFuzz))
		t.Targets = append(t.Targets, reqsByFuzz...)
	}

	var initTasks []*model.Request
	for _, req := range t.Targets {
		// 保存所有任务结果
		t.SaveAllReqInfo(req)

		if t.smartFilter.DoFilter(req) {
			logger.Logger.Debugf("filter req: " + req.URL.RequestURI())
			continue
		}
		initTasks = append(initTasks, req)
		// 保存有效的请求任务信息
		t.SaveReqResult(req)
	}
	logger.Logger.Info("filter repeat, target count: ", len(initTasks))

	for _, req := range initTasks {
		if !engine2.IsIgnoredByKeywordMatch(*req, t.Config.IgnoreKeywords) {
			t.addTask2Pool(req)
		}
	}

	t.taskWG.Wait()
	t.Result.Close()
	// for index := range t.Result.AllReqList {
	// 	todoFilterAll[index] = t.Result.AllReqList[index]
	// }

}

/**
添加任务到协程池
添加之前实时过滤
*/
func (t *CrawlerTask) addTask2Pool(req *model.Request) {
	t.taskCountLock.Lock()
	if t.crawledCount >= t.Config.MaxCrawlCount {
		t.taskCountLock.Unlock()
		return
	} else {
		t.crawledCount += 1
	}
	t.taskCountLock.Unlock()

	t.taskWG.Add(1)
	task := t.generateTabTask(req)
	go func() {
		err := t.Pool.Submit(task.Task)
		if err != nil {
			t.taskWG.Done()
			logger.Logger.Error("addTask2Pool ", err)
		}
	}()
}

/**
单个运行的tab标签任务，实现了workpool的接口
*/
func (t *tabTask) Task() {
	defer t.crawlerTask.taskWG.Done()
	tab := engine2.NewTab(t.browser, *t.req, engine2.TabConfig{
		TabRunTimeout:           t.crawlerTask.Config.TabRunTimeout,
		DomContentLoadedTimeout: t.crawlerTask.Config.DomContentLoadedTimeout,
		EventTriggerMode:        t.crawlerTask.Config.EventTriggerMode,
		EventTriggerInterval:    t.crawlerTask.Config.EventTriggerInterval,
		BeforeExitDelay:         t.crawlerTask.Config.BeforeExitDelay,
		EncodeURLWithCharset:    t.crawlerTask.Config.EncodeURLWithCharset,
		IgnoreKeywords:          t.crawlerTask.Config.IgnoreKeywords,
		CustomFormValues:        t.crawlerTask.Config.CustomFormValues,
		CustomFormKeywordValues: t.crawlerTask.Config.CustomFormKeywordValues,
	})
	tab.Start()

	// 收集结果
	t.crawlerTask.Result.resultLock.Lock()
	// 保存所有结果，包括域名
	for _, req := range tab.ResultList {
		t.crawlerTask.SaveAllReqInfo(req)
	}
	t.crawlerTask.Result.resultLock.Unlock()

	for _, req := range tab.ResultList {
		if t.crawlerTask.Config.FilterMode == config.SimpleFilterMode {
			if !t.crawlerTask.smartFilter.SimpleFilter.DoFilter(req) {
				t.crawlerTask.Result.resultLock.Lock()
				// 保存有效请求结果
				t.crawlerTask.SaveReqResult(req)
				t.crawlerTask.Result.resultLock.Unlock()
				if !engine2.IsIgnoredByKeywordMatch(*req, t.crawlerTask.Config.IgnoreKeywords) {
					t.crawlerTask.addTask2Pool(req)
				}
			}
		} else {
			if !t.crawlerTask.smartFilter.DoFilter(req) {
				t.crawlerTask.Result.resultLock.Lock()
				// 保存有效请求结果
				t.crawlerTask.SaveReqResult(req)
				t.crawlerTask.Result.resultLock.Unlock()
				if !engine2.IsIgnoredByKeywordMatch(*req, t.crawlerTask.Config.IgnoreKeywords) {
					t.crawlerTask.addTask2Pool(req)
				}
			}
		}
	}
	tab.ResultList = nil
}

func getOption(postData string, taskConfig taskPkg.TaskConfig) model.Options {
	var option model.Options
	if postData != "" {
		option.PostData = postData
	}
	if taskConfig.ExtraHeadersString != "" {
		err := json.Unmarshal([]byte(taskConfig.ExtraHeadersString), &taskConfig.ExtraHeaders)
		if err != nil {
			logger.Logger.Fatal("custom headers can't be Unmarshal.")
			panic(err)
		}
		option.Headers = taskConfig.ExtraHeaders
	}
	return option
}

// newReqForUrl 根据 url 和配置构造任务信息
func newReqForUrl(targetUrl string, taskConf taskPkg.TaskConfig, postData string) *model.Request {
	var req model.Request
	url, err := model.GetUrl(targetUrl)
	if err != nil {
		logger.Logger.Error("parse url failed, ", err)
		return nil
	}
	if postData != "" {
		req = model.GetRequest(config.POST, url, getOption(postData, taskConf))
	} else {
		req = model.GetRequest(config.GET, url, getOption(postData, taskConf))
	}
	req.Proxy = taskConf.Proxy
	return &req
}

func MakeTargets(urls []string, taskConf taskPkg.TaskConfig, postData string) []*model.Request {
	var targets []*model.Request
	for _, _url := range urls {
		req := newReqForUrl(_url, taskConf, postData)
		targets = append(targets, req)
	}
	return targets
}

// GetRedisTarget 从 redis 中获取任务信息
func GetRedisTarget(taskConf taskPkg.TaskConfig, postData string) ([]*model.Request, *biz.EngineUsecase) {
	var targets []*model.Request
	enguc := wireRedisCacheRepo(taskConf)
	tasks := enguc.GetTasks()
	for _, t := range tasks {
		req := newReqForUrl(t.Url, taskConf, postData)
		// TODO: set id
		req.TaskID = t.ID
		targets = append(targets, req)
	}
	return targets, enguc
}

// SaveAllReqInfo 保存域名结果、爬虫请求结果
func (t *CrawlerTask) SaveAllReqInfo(req *model.Request) {
	var wg sync.WaitGroup
	wg.Add(3)
	// 结果保存
	go func() {
		// 保存所有结果
		defer wg.Done()
		if t.Result.AllReqSave == nil || t.Result.AllReqSimpFilter.DoFilter(req) {
			return
		}
		t.Result.AllReqSave.Save(req)
	}()

	// 全部域名保存
	go func() {
		defer wg.Done()
		if t.Result.allDomainSave != nil {
			t.Result.allDomainSave.Save(req)
		}
	}()

	// 子域名保存
	go func() {
		defer wg.Done()
		if t.Result.subDomainSave != nil {
			t.Result.subDomainSave.Save(req)
		}
	}()
	wg.Wait()
}

// SaveReqResult 保存有效的请求结果
func (t *CrawlerTask) SaveReqResult(req *model.Request) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if t.Result.ReqSave != nil {
			t.Result.ReqSave.Save(req)
		}
	}()

	go func() {
		defer wg.Done()
		if t.redisUsecase != nil {
			t.redisUsecase.SetTaskResult(req)
		}
	}()
	wg.Wait()
}

func (r *Result) Close() {
	r.resultLock.Lock()
	defer r.resultLock.Unlock()
	if r.allDomainSave != nil {
		r.allDomainSave.Close()
	}
	if r.subDomainSave != nil {
		r.subDomainSave.Close()
	}
	if r.ReqSave != nil {
		r.ReqSave.Close()
	}
	if r.AllReqSave != nil {
		r.AllReqSave.Close()
	}
}
