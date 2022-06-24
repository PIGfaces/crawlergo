package engine_test

import (
	"context"
	"log"
	"path"
	"strings"
	"testing"

	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

var (
	chromiumPath = ""
	noHeadless   = true
)

func TestTab(t *testing.T) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],

		// 执行路径
		chromedp.ExecPath(chromiumPath),
		// 无头模式
		chromedp.Flag("headless", !noHeadless),
		// 禁用GPU，不显示GUI
		chromedp.Flag("disable-gpu", true),
		// 隐身模式启动
		chromedp.Flag("incognito", true),
		// 取消沙盒模式
		chromedp.Flag("no-sandbox", true),
		// 忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),

		chromedp.Flag("disable-images", true),
		//
		chromedp.Flag("disable-web-security", true),
		//
		chromedp.Flag("disable-xss-auditor", true),
		//
		chromedp.Flag("disable-setuid-sandbox", true),

		chromedp.Flag("allow-running-insecure-content", true),

		chromedp.Flag("disable-webgl", true),

		chromedp.Flag("disable-popup-blocking", true),

		chromedp.WindowSize(1920, 1080),
	)
	// 设置浏览器代理

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	bctx, _ := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	// https://github.com/chromedp/chromedp/issues/824#issuecomment-845664441
	// 如果需要在一个浏览器上创建多个tab，则需要先创建浏览器的上下文，即运行下面的语句
	err := chromedp.Run(bctx)
	assert.Nil(t, err)

}

func TestStaticUrl(t *testing.T) {
	url := "http://www.baidu.com/test/test.png"
	_UrlMod, err := model.GetUrl(url)
	assert.Nil(t, err)
	assert.Equal(t, ".png", path.Ext(_UrlMod.Path))
	assert.Equal(t, "test.png", path.Base(_UrlMod.Path))
}

func TestFilterStatic(t *testing.T) {
	url := "http://www.baidu.com/test/Modifypng"
	_UrlMod, err := model.GetUrl(url)
	assert.Nil(t, err)
	assert.Equal(t, "", path.Ext(_UrlMod.Path))
	for _, suffix := range config.StaticSuffix {
		assert.Equal(t, false, strings.HasSuffix(strings.ToLower(_UrlMod.Path), suffix))
	}
}
