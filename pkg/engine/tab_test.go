package engine_test

import (
	"context"
	"log"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/engine"
	"github.com/PIGfaces/crawlergo/pkg/logger"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/chromedp/cdproto/fetch"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	chromiumPath = ""
	// chromiumPath = "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	noHeadless = true
)

func getBrowserCtx(t *testing.T) (context.Context, context.CancelFunc, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],

		// 执行路径
		chromedp.ExecPath(chromiumPath),
		// 无头模式
		chromedp.Flag("headless", !noHeadless),
	)
	// 设置浏览器代理

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()
	bctx, bcancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	// https://github.com/chromedp/chromedp/issues/824#issuecomment-845664441
	// 如果需要在一个浏览器上创建多个tab，则需要先创建浏览器的上下文，即运行下面的语句
	err := chromedp.Run(bctx)
	assert.Nil(t, err)
	return bctx, bcancel, cancel
}

func getTabTask(t *testing.T, ctx context.Context, cancel context.CancelFunc, url string) engine.Tab {
	urlObj, err := model.GetUrl(url)
	assert.Nil(t, err)
	req := model.GetRequest(http.MethodGet, urlObj)

	tab := engine.Tab{
		Ctx:         ctx,
		Cancel:      cancel,
		NavigateReq: &req,
	}
	logger.Logger.SetLevel(logrus.DebugLevel)
	return tab
}

// func TestEventResponseReceived() {

// }

func TestEventResponseReceived(t *testing.T) {
	bctx, bcancel, cancel := getBrowserCtx(t)
	defer cancel()
	tab := getTabTask(t, bctx, bcancel, "http://www.caih.com")

	chromedp.ListenTarget(bctx, func(ev interface{}) {
		switch v := ev.(type) {
		case *network.EventResponseReceived:
			// t.Log("mimeType: ", v.Response.MimeType, " status: ", v.Response.Status, " parse url: ", v.Response.URL)
			if v.Response.MimeType == "application/javascript" || v.Response.MimeType == "text/html" || v.Response.MimeType == "application/json" {
				tab.WG.Add(1)
				go tab.ParseResponseURL(v)
			}
			if v.RequestID.String() == tab.NavNetworkID {
				tab.WG.Add(1)
				go tab.GetContentCharset(v)
			}
		}
	})

	err := chromedp.Run(bctx, chromedp.Navigate("https://www.caih.com"))
	assert.Nil(t, err)
	t.Log("result len: ", len(tab.ResultList), " ", tab.ResultList)
}

func TestEventAuthRequired(t *testing.T) {
	bctx, bcancel, cancel := getBrowserCtx(t)
	defer cancel()
	tab := getTabTask(t, bctx, bcancel, "http://www.caih.com")
	chromedp.ListenTarget(bctx, func(ev interface{}) {
		switch v := ev.(type) {
		case *fetch.EventAuthRequired:
			t.Log("get auth required")
			// t.Log("mimeType: ", v.Response.MimeType, " status: ", v.Response.Status, " parse url: ", v.Response.URL)
			tab.HandleAuthRequired(v)
		}
	})

	err := chromedp.Run(bctx, chromedp.Navigate(tab.NavigateReq.URL.String()))
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
