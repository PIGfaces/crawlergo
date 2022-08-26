package engine

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/PIGfaces/crawlergo/pkg/logger"

	"github.com/chromedp/cdproto/browser"
	"github.com/chromedp/chromedp"
)

type Browser struct {
	Ctx          context.Context
	Cancel       context.CancelFunc
	tabs         []context.Context
	tabCancels   []context.CancelFunc
	ExtraHeaders map[string]interface{}
	lock         sync.Mutex
}

func init() {

}

func InitBrowser(chromiumPath string, incognito bool, extraHeaders map[string]interface{}, proxy string, noHeadless bool) *Browser {
	var bro Browser
	opts := append(chromedp.DefaultExecAllocatorOptions[:],

		// 执行路径
		// 无头模式
		chromedp.Flag("headless", !noHeadless),
		// 禁用GPU，不显示GUI
		chromedp.Flag("disable-gpu", true),
		// 隐身模式启动
		chromedp.Flag("incognito", incognito),
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
	if chromiumPath != "" {
		opts = append(opts, chromedp.ExecPath(chromiumPath))
	}
	if proxy != "" {
		opts = append(opts, chromedp.ProxyServer(proxy))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	bctx, _ := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(log.Printf),
		// chromedp.WithDebugf(log.Printf),
	)
	// https://github.com/chromedp/chromedp/issues/824#issuecomment-845664441
	// 如果需要在一个浏览器上创建多个tab，则需要先创建浏览器的上下文，即运行下面的语句
	if err := chromedp.Run(bctx); err != nil {
		logger.Logger.Fatal("open browser failed: [", chromiumPath, "] error: ", err)
	}
	bro.Cancel = cancel
	bro.Ctx = bctx
	bro.ExtraHeaders = extraHeaders
	return &bro
}

func (bro *Browser) NewTab(timeout time.Duration) (context.Context, context.CancelFunc) {
	bro.lock.Lock()
	defer bro.lock.Unlock()
	ctx, cancel := chromedp.NewContext(bro.Ctx)
	//defer cancel()
	tCtx, _ := context.WithTimeout(ctx, timeout)
	bro.tabs = append(bro.tabs, tCtx)
	bro.tabCancels = append(bro.tabCancels, cancel)
	//defer cancel2()

	//return bro.Ctx, &cancel
	return tCtx, cancel
}

func (bro *Browser) Close() {
	logger.Logger.Info("closing browser.")
	for _, cancel := range bro.tabCancels {
		cancel()
	}

	for _, ctx := range bro.tabs {
		err := browser.Close().Do(ctx)
		if err != nil {
			logger.Logger.Debug(err)
		}
	}

	err := browser.Close().Do(bro.Ctx)
	if err != nil {
		logger.Logger.Debug(err)
	}
	bro.Cancel()
}
