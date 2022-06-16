package main

//go:generate mockgen -source=main.go -destination=../../mock/cmd/crawlergo/main_mock.go -package=main
import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"runtime/pprof"

	"github.com/PIGfaces/crawlergo/pkg"
	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/logger"
	model2 "github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/PIGfaces/crawlergo/pkg/tools"
	"github.com/PIGfaces/crawlergo/pkg/tools/requests"
	"github.com/shirou/gopsutil/process"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	taskPkg "github.com/PIGfaces/crawlergo/pkg/task"
)

/**
命令行调用适配器

用于生成开源的二进制程序
*/

type Result struct {
	ReqList       []Request `json:"req_list"`
	AllReqList    []Request `json:"all_req_list"`
	AllDomainList []string  `json:"all_domain_list"`
	SubDomainList []string  `json:"sub_domain_list"`
}

type Request struct {
	Url     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]interface{} `json:"headers"`
	Data    string                 `json:"data"`
	Source  string                 `json:"source"`
}

type ProxyTask struct {
	req       *model2.Request
	pushProxy string
}

type ProcessImp interface {
	CPUPercent() (float64, error)
	MemoryPercent() (float32, error)
}

type scaleWeight struct {
	CpuWeight    float64 `json:"CpuWeight"`
	MemWeight    float64 `json:"MemWeight"`
	TabTTLWeight float64 `json:"TabTTLWeight"`
}

const DefaultMaxPushProxyPoolMax = 10
const DefaultLogLevel = "Info"

var (
	taskConfig              taskPkg.TaskConfig
	outputMode              string
	postData                string
	signalChan              chan os.Signal
	ignoreKeywords          *cli.StringSlice
	customFormTypeValues    *cli.StringSlice
	customFormKeywordValues *cli.StringSlice
	pushAddress             string
	pushProxyPoolMax        int
	pushProxyWG             sync.WaitGroup
	logLevel                string
	Version                 string
	isCPUPprof              bool
	isMemPprof              bool
	autoScaleTabs           bool
	scaleCtx                context.Context
	scaleCtxCancle          context.CancelFunc
	sweight                 string
	sw                      scaleWeight = scaleWeight{
		CpuWeight:    0.4,
		MemWeight:    0.3,
		TabTTLWeight: 0.3,
	}
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Logger.Error("occur some panic: ", r)
			// 要释放资源，正常退出
			signalChan <- syscall.SIGTERM
		}
	}()
	author := cli.Author{
		Name:  "9ian1i, chongzhenghuang",
		Email: "9ian1itp@gmail.com, chongzhenghuang@caih.com",
	}

	ignoreKeywords = cli.NewStringSlice(config.DefaultIgnoreKeywords...)
	customFormTypeValues = cli.NewStringSlice()
	customFormKeywordValues = cli.NewStringSlice()

	// app := cli.NewApp()

	app := &cli.App{
		Name:      "crawlergo",
		Usage:     "A powerful browser crawler for web vulnerability scanners",
		UsageText: "crawlergo [global options] url1 url2 url3 ... (must be same host)",
		Version:   Version,
		Authors:   []*cli.Author{&author},
		Flags:     cliFlags,
		Action:    run,
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Logger.Fatal(err)
	}
}

func run(c *cli.Context) error {
	saveCpuPprof()
	saveMemPprof()

	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	scaleCtx, scaleCtxCancle = context.WithCancel(context.Background())

	if c.Args().Len() == 0 && taskConfig.RedisConnectInfo == "" {
		// if c.Args().Len() == 0 {
		return fmt.Errorf("url must be set: %d", c.Args().Len())
	}

	// 设置日志输出级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Logger.Fatal(err)
	}
	logger.Logger.SetLevel(level)

	taskConfig.IgnoreKeywords = ignoreKeywords.Value()
	if taskConfig.Proxy != "" {
		logger.Logger.Info("request with proxy: ", taskConfig.Proxy)
	}

	// 检查自定义的表单参数配置
	taskConfig.CustomFormValues, err = parseCustomFormValues(customFormTypeValues.Value())
	if err != nil {
		logger.Logger.Fatal(err)
	}
	taskConfig.CustomFormKeywordValues, err = keywordStringToMap(customFormKeywordValues.Value())
	if err != nil {
		logger.Logger.Fatal(err)
	}

	// 开始爬虫任务
	task, err := pkg.NewCrawlerTask(c.Args().Slice(), taskConfig, postData)
	if err != nil {
		logger.Logger.Error("create crawler task failed.")
		os.Exit(-1)
	}

	// 提示自定义表单填充参数
	if len(taskConfig.CustomFormValues) > 0 {
		logger.Logger.Info("Custom form values, " + tools.MapStringFormat(taskConfig.CustomFormValues))
	}
	// 提示自定义表单填充参数
	if len(taskConfig.CustomFormKeywordValues) > 0 {
		logger.Logger.Info("Custom form keyword values, " + tools.MapStringFormat(taskConfig.CustomFormKeywordValues))
	}
	if _, ok := taskConfig.CustomFormValues["default"]; !ok {
		logger.Logger.Info("If no matches, default form input text: " + config.DefaultInputText)
		taskConfig.CustomFormValues["default"] = config.DefaultInputText
	}

	go handleExit(task)
	// ModifyScaleWeight 调整 scale 权重
	ModifyScaleWeight(sweight)
	if autoScaleTabs || sweight != "" {
		go autoScaleConcurrency(task)
	}
	logger.Logger.Info("Start crawling.")
	task.Run()
	// result := task.Result

	// 内置请求代理
	// if pushAddress != "" {
	// 	logger.Logger.Info("pushing results to ", pushAddress, ", max pool number:", pushProxyPoolMax)
	// 	Push2Proxy(result.ReqList)
	// }

	return nil
}

func parseCustomFormValues(customData []string) (map[string]string, error) {
	parsedData := map[string]string{}
	for _, item := range customData {
		keyValue := strings.Split(item, "=")
		if len(keyValue) < 2 {
			return nil, errors.New("invalid form item: " + item)
		}
		key := keyValue[0]
		if !tools.StringSliceContain(config.AllowedFormName, key) {
			return nil, errors.New("not allowed form key: " + key)
		}
		value := keyValue[1]
		parsedData[key] = value
	}
	return parsedData, nil
}

func keywordStringToMap(data []string) (map[string]string, error) {
	parsedData := map[string]string{}
	for _, item := range data {
		keyValue := strings.Split(item, "=")
		if len(keyValue) < 2 {
			return nil, errors.New("invalid keyword format: " + item)
		}
		key := keyValue[0]
		value := keyValue[1]
		parsedData[key] = value
	}
	return parsedData, nil
}

/**
原生被动代理推送支持
*/
func Push2Proxy(reqList []*model2.Request) {
	pool, _ := ants.NewPool(pushProxyPoolMax)
	defer pool.Release()
	for _, req := range reqList {
		task := ProxyTask{
			req:       req,
			pushProxy: pushAddress,
		}
		pushProxyWG.Add(1)
		go func() {
			err := pool.Submit(task.doRequest)
			if err != nil {
				logger.Logger.Error("add Push2Proxy task failed: ", err)
				pushProxyWG.Done()
			}
		}()
	}
	pushProxyWG.Wait()
}

/**
协程池请求的任务
*/
func (p *ProxyTask) doRequest() {
	defer pushProxyWG.Done()
	_, _ = requests.Request(p.req.Method, p.req.URL.String(), tools.ConvertHeaders(p.req.Headers), []byte(p.req.PostData),
		&requests.ReqOptions{Timeout: 1, AllowRedirect: false, Proxy: p.pushProxy})
}

func autoScaleConcurrency(t *pkg.CrawlerTask) {
	tc := time.NewTicker(time.Second * 10)
	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		logger.Logger.Fatal("cannot start monitor process self: ", err.Error())
	}
	for {
		select {
		case <-tc.C:
			ModifyPool(t, p)
		case <-scaleCtx.Done():
			logger.Logger.Debug("scale tab max number exit...")
			tc.Stop()
			return
		}
	}
}

// 自动调整协程池
func ModifyPool(t *pkg.CrawlerTask, p ProcessImp) {
	cpuPer, err := p.CPUPercent()
	if err != nil {
		logger.Logger.Debug("get cpu percent failed: ", err.Error())
	}
	memPer, err := p.MemoryPercent()
	if err != nil {
		logger.Logger.Debug("get mem percent failed: ", err.Error())
	}
	// 标签页超时比例
	tabTimeoutPercent := float64(t.TabTTLNum) / float64(t.TabNum)

	// 计算权重: cpu 利用率占 40%, 内存利用率占 30%, 标签页超时占 30%
	totalPer := (cpuPer*sw.CpuWeight + float64(memPer)*sw.MemWeight + tabTimeoutPercent*sw.TabTTLWeight)
	// 每次伸缩都以 30 个标签页为基准
	scaleNum := math.Ceil(30 * totalPer)
	needModifyTabNumFlag := false
	if t.Pool.Running() == t.Config.MaxTabsCount {
		// 若所有的标签页都跑满才会进行伸缩
		if totalPer > 0.7 {
			// 若负载权重 > 70% 说明机器负载比较大, 降低
			if int(scaleNum) < t.Config.MaxTabsCount {
				// 正常伸缩
				needModifyTabNumFlag = true
				t.Config.MaxTabsCount -= int(scaleNum)
			} else {
				// 折半
				needModifyTabNumFlag = true
				t.Config.MaxTabsCount -= int(math.Ceil(float64(t.Config.MaxTabsCount) / 2))
			}
		} else if totalPer < 0.4 {
			needModifyTabNumFlag = true
			t.Config.MaxTabsCount += int(scaleNum)
		}
	}

	if needModifyTabNumFlag {
		// 对并发进行调整
		t.Pool.Tune(t.Config.MaxTabsCount)
		logger.Logger.Debug("adjuest tab runing num: ", t.Config.MaxTabsCount)
		return
	}
	logger.Logger.Debug("no adjuest tab num")
}

func handleExit(t *pkg.CrawlerTask) {
	<-signalChan
	scaleCtxCancle()
	logger.Logger.Debug("exit ...")
	t.Pool.Tune(1)
	t.Pool.Release()
	t.Browser.Close()
	t.Result.Close()
	os.Exit(-1)
}

// saveCpuPprof 保存 cpu 的信息
func saveCpuPprof() {

	if isCPUPprof {
		file, err := os.Create(fmt.Sprintf("./%s_cpu.pprof", time.Now().String()))
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(file)
		defer func() {
			pprof.StopCPUProfile()
			file.Close()
		}()
	}
}

// saveMemPprof 保存内存信息
func saveMemPprof() {
	if isMemPprof {
		file, err := os.Create(fmt.Sprintf("./%s_mem.pprof", time.Now().Format("2022-06-13_13:00:00")))
		if err != nil {
			panic(err)
		}
		pprof.WriteHeapProfile(file)
		defer file.Close()
	}
}

func ModifyScaleWeight(weightJson string) scaleWeight {
	if len(weightJson) <= 0 {
		return scaleWeight{}
	}
	destWeight := scaleWeight{}
	if err := json.Unmarshal([]byte(weightJson), &destWeight); err != nil {
		logger.Logger.Fatal("scale weight cannot unmarshal, please check your json struct: ", err.Error())
	}

	for _, v := range []float64{
		destWeight.CpuWeight,
		destWeight.MemWeight,
		destWeight.TabTTLWeight,
	} {
		if v > 0.9 || v < 0.1 {
			logger.Logger.Fatal("signal weight must in (0.1, 0.9)")
		}
	}
	sw = destWeight
	return destWeight
}
