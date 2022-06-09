package filter

import (
	"testing"

	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/model"

	"github.com/stretchr/testify/assert"
)

var (
	baseUri = "http://testhtml5.vulnweb.com/#/popular/page/-2"
	sameUri = "http://testhtml5.vulnweb.com/#/popular/page/0"

	fragmentSimpleUri   = "http://test.com/#/search?keyword=Crawlergo&source=1"
	fragmentMultWordUri = "http://test.com/#/search?keyword=CrawlergoCrawlergoCrawlergo&source=1"

	smart = SmartFilter{
		SimpleFilter: SimpleFilter{
			HostLimit: "testhtml5.vulnweb.com",
		},
	}
)

func TestMarkPath(t *testing.T) {
	// 测试实例

	baseUrlMarkResult := MarkPath(baseUri)
	sameUrlMarkResult := MarkPath(sameUri)

	t.Log("mark path result: ", baseUrlMarkResult)
	assert.Equal(t, baseUrlMarkResult, sameUrlMarkResult)
}

func TestMarkPathForFragment(t *testing.T) {
	baseReq := getRequest(t, baseUri)
	sameReq := getRequest(t, sameUri)

	fragmentSimpleReq := getRequest(t, fragmentSimpleUri)
	fragmentMultWordReq := getRequest(t, fragmentMultWordUri)

	t.Log("\nfragment URL: ", baseReq.URL.Fragment, "\tMarked: ", MarkPath(baseReq.URL.Fragment))
	t.Log("\nsame multi word at fragment: ", fragmentSimpleReq.URL.Fragment, "\tMarked: ", MarkPath(fragmentSimpleReq.URL.Fragment))
	assert.Equal(t, MarkPath(baseReq.URL.Fragment), MarkPath(sameReq.URL.Fragment))
	assert.Equal(t, MarkPath(fragmentSimpleReq.URL.Fragment), MarkPath(fragmentMultWordReq.URL.Fragment))
}

func TestDoFilter(t *testing.T) {
	// ==========> init data
	baseUrl, err := model.GetUrl(baseUri)
	assert.Nil(t, err)

	sameUrl, err := model.GetUrl(sameUri)
	assert.Nil(t, err)

	baseReq := model.GetRequest("get", baseUrl, model.Options{
		Headers: map[string]interface{}{
			"Referer":     "http://testhtml5.vulnweb.com/#/popular/page/-1",
			"Spider-Name": "crawlergo",
			"User-Agent":  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36",
		},
	})

	sameReq := model.GetRequest("get", sameUrl, model.Options{
		Headers: map[string]interface{}{
			"Referer":     "http://testhtml5.vulnweb.com/#/popular/page/-1",
			"Spider-Name": "crawlergo",
			"User-Agent":  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36",
		},
	})

	baseReq.Source = config.FromDOM
	sameReq.Source = config.FromDOM
	// filter create

	smart.Init()
	// <========
	// 计算了两次
	assert.Equal(t, smart.DoFilter(&baseReq), false)
	// 计算了一次
	assert.Equal(t, smart.DoFilter(&sameReq), true)

	// assert.Equal(t, baseReq.UniqueId(), sameReq.UniqueId())
}

func TestGetPathID(t *testing.T) {
	baseUrl, err := model.GetUrl(baseUri)
	assert.Nil(t, err)

	sameUrl, err := model.GetUrl(sameUri)
	assert.Nil(t, err)

	assert.Equal(t, getPathID(baseUrl.Path), getPathID(sameUrl.Path))
}

func TestGetMark(t *testing.T) {
	baseReq := getRequest(t, baseUri)

	sameReq := getRequest(t, sameUri)

	smart.getMark(&baseReq)
	smart.getMark(&sameReq)

	assert.Equal(t, baseReq.Filter.UniqueId, sameReq.Filter.UniqueId)
}

func TestRepeatCountStatistic(t *testing.T) {
	baseReq := getRequest(t, baseUri)

	sameReq := getRequest(t, sameUri)

	smart.getMark(&baseReq)
	smart.getMark(&sameReq)
	smart.repeatCountStatistic(&baseReq)
	smart.repeatCountStatistic(&baseReq)
}

// TestGetAllMark 测试 get 方法所有打上标记
func TestGetAllMarkAndCalcUniqueID(t *testing.T) {
	baseReq := getRequest(t, baseUri)

	sameReq := getRequest(t, sameUri)

	smart.getMark(&baseReq)
	smart.getMark(&sameReq)
	smart.repeatCountStatistic(&baseReq)
	smart.repeatCountStatistic(&sameReq)

	t.Log(baseReq.Filter.UniqueId)
	assert.Equal(t, baseReq.Filter.UniqueId, sameReq.Filter.UniqueId)
}

func getRequest(t *testing.T, url string) model.Request {
	baseUrl, err := model.GetUrl(url)
	assert.Nil(t, err)

	baseReq := model.GetRequest("get", baseUrl, model.Options{
		Headers: map[string]interface{}{
			"Referer":     "http://testhtml5.vulnweb.com/#/popular/page/-1",
			"Spider-Name": "crawlergo",
			"User-Agent":  "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.0 Safari/537.36",
		},
	})

	return baseReq
}
