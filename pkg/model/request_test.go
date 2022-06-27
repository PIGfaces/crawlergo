package model_test

import (
	"testing"

	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	url := "http://www.baidu.com"
	UrlObj, err := model.GetUrl(url)
	assert.Nil(t, err)

	reqObj := model.GetRequest(config.GET, UrlObj)
	reqObj.Headers["test"] = "test"

	cpReq := reqObj.Copy()
	assert.Equal(t, reqObj.Headers["test"], cpReq.Headers["test"])
}
