package engine

import (
	"time"

	"github.com/PIGfaces/crawlergo/pkg/config"
	taskConf "github.com/PIGfaces/crawlergo/pkg/task"
)

type TabConfig struct {
	TabRunTimeout           time.Duration
	DomContentLoadedTimeout time.Duration
	EventTriggerMode        string        // 事件触发的调用方式： 异步 或 顺序
	EventTriggerInterval    time.Duration // 事件触发的间隔 单位毫秒
	BeforeExitDelay         time.Duration // 退出前的等待时间，等待DOM渲染，等待XHR发出捕获
	EncodeURLWithCharset    bool
	SaveHtmlCode            bool
	IgnoreKeywords          []string //
	Proxy                   string
	CustomFormValues        map[string]string
	CustomFormKeywordValues map[string]string
	UploadFiles             []string
}

type TabConfigOptFunc func(*TabConfig)

func NewByCrawTab(tconf *taskConf.TaskConfig) *TabConfig {
	return NewTabConfig(
		WithTabRunTimeout(tconf.TabRunTimeout),
		WithDomContentLoadedTimeout(tconf.DomContentLoadedTimeout),
		WithEventTriggerInterval(tconf.EventTriggerInterval),
		WithEventTriggerMode(tconf.EventTriggerMode),
		WithBeforeExitDelay(tconf.BeforeExitDelay),
		WithEncodeURLWithCharset(tconf.EncodeURLWithCharset),
		WithIgnoreKeywords(tconf.IgnoreKeywords),
		WithCustomFormValues(tconf.CustomFormValues),
		WithCustomFormKeywordValues(tconf.CustomFormKeywordValues),
		WithSaveHtmlCode(tconf.RedisConnectInfo),
	)
}

func NewTabConfig(optFunc ...TabConfigOptFunc) *TabConfig {
	var conf TabConfig
	for _, fn := range optFunc {
		fn(&conf)
	}
	return &conf
}

func WithProxy(proxy string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.Proxy = proxy
	}
}

func WithEventTriggerMode(mod string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.EventTriggerMode = mod
	}
}

func WithCustomFormValues(customFormValue map[string]string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.CustomFormValues = customFormValue
	}
}

func WithCustomFormKeywordValues(customKeyFromKey map[string]string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.CustomFormKeywordValues = customKeyFromKey
	}
}

func WithIgnoreKeywords(ignoreKeywords []string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.IgnoreKeywords = ignoreKeywords
	}
}

func WithEncodeURLWithCharset(isEncode bool) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.EncodeURLWithCharset = isEncode
	}
}

func WithTabRunTimeout(duration time.Duration) TabConfigOptFunc {
	return func(tc *TabConfig) {
		if duration != 0 {
			tc.TabRunTimeout = duration
		} else {
			tc.TabRunTimeout = config.TabRunTimeout
		}
	}
}
func WithDomContentLoadedTimeout(duration time.Duration) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.DomContentLoadedTimeout = duration
	}
}
func WithEventTriggerInterval(duration time.Duration) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.EventTriggerInterval = duration
	}
}
func WithBeforeExitDelay(duration time.Duration) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.BeforeExitDelay = duration
	}
}

func WithSaveHtmlCode(redisConnInfo string) TabConfigOptFunc {
	return func(tc *TabConfig) {
		tc.SaveHtmlCode = len(redisConnInfo) > 0
	}
}
