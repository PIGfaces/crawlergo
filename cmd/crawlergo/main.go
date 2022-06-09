package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/PIGfaces/crawlergo/pkg"
	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/logger"
	model2 "github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/PIGfaces/crawlergo/pkg/tools"
	"github.com/PIGfaces/crawlergo/pkg/tools/requests"

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
)

func main() {
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
	signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

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

func handleExit(t *pkg.CrawlerTask) {
	<-signalChan
	fmt.Println("exit ...")
	t.Pool.Tune(1)
	t.Pool.Release()
	t.Browser.Close()
	t.Result.Close()
	os.Exit(-1)
}
