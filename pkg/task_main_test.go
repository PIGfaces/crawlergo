package pkg

import (
	"os"
	"strings"
	"testing"

	"github.com/PIGfaces/crawlergo/pkg/config"
	"github.com/PIGfaces/crawlergo/pkg/model"
	"github.com/stretchr/testify/assert"
)

var (
	fileNames = []string{"image.png", "demo.txt", "hello.mp3"}
	dirNames  = []string{"demo", "test", "hello"}
)

func getRequest(t *testing.T, urls ...string) []*model.Request {
	reqList := []*model.Request{}
	for _, item := range urls {
		url, err := model.GetUrl(item)
		assert.Nil(t, err)
		req := model.GetRequest(config.GET, url)
		reqList = append(reqList, &req)
	}
	return reqList
}

func TestSetUploadFile(t *testing.T) {
	testDir := "./testDir"
	pwd, err := os.Getwd()
	assert.Nil(t, err)

	err = os.Mkdir(pwd+"/testDir", 0777)
	assert.Nil(t, err)
	defer os.RemoveAll(testDir)

	// 创建测试文件
	var testFileListPath = []string{}
	for _, name := range fileNames {
		var testPath = pwd + "/testDir/" + name
		f, err := os.OpenFile(testPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
		// _, err = os.Create(testPath)
		f.Close()
		assert.Nil(t, err, err)
		testFileListPath = append(testFileListPath, testPath)
	}

	// 创建测试目录
	for _, name := range dirNames {
		err = os.Mkdir(pwd+"/testDir/"+name, 0777)
		assert.Nil(t, err, err)
	}

	ct := CrawlerTask{}

	// 相对路径测试
	ct.setUploadFileDir(testDir)
	assert.Equal(t, len(ct.UploadFiles), len(testFileListPath))

	var set = make(map[string]struct{})
	for _, val := range testFileListPath {
		set[val] = struct{}{}
	}
	// result should equal
	for _, val := range ct.UploadFiles {
		_, ok := set[val]
		assert.Equal(t, ok, true)
	}

	// 完整路径测试
	ct.setUploadFileDir(pwd + "/testDir")
	assert.Equal(t, len(ct.UploadFiles), len(testFileListPath))

	for _, val := range testFileListPath {
		set[val] = struct{}{}
	}
	// result should equal
	for _, val := range ct.UploadFiles {
		_, ok := set[val]
		assert.Equal(t, ok, true)
	}

	// 文件测试
	ct.setUploadFileDir(testFileListPath[0])

	assert.Equal(t, ct.UploadFiles, []string{testFileListPath[0]})
}

func TestXxx(t *testing.T) {
	path := "/test/demo/hello/world"
	t.Log(strings.Split(path, "/"))
}

func TestFuzzTabResultList(t *testing.T) {
	testReqList := getRequest(t, "http://testphp.vulnweb.com/Mod_Rewrite_Shop/Details/network-attached-storage-dlink/1/", "http://testphp.vulnweb.com/Mod_Rewrite_Shop/Details/network-attached-storage-dlink/2/")
	resultList := fuzzTabResultList(testReqList)
	for _, v := range resultList {
		t.Log(*v)
	}
	assert.Len(t, resultList, 4)
}
