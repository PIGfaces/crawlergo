package resultsave

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/PIGfaces/crawlergo/pkg/logger"
	"github.com/PIGfaces/crawlergo/pkg/model"
)

type (
	FileSave struct {
		file *os.File
		iow  *bufio.Writer
	}

	RequestInfo struct {
		Url     string                 `json:"url"`
		Method  string                 `json:"method"`
		Headers map[string]interface{} `json:"headers"`
		Data    string                 `json:"data"`
		Source  string                 `json:"source"`
	}
)

func getJsonSerialize(req model.Request) ([]byte, error) {
	result := RequestInfo{
		Method:  req.Method,
		Url:     req.URL.String(),
		Data:    req.PostData,
		Headers: req.Headers,
		Source:  req.Source,
	}
	resBytes, err := json.Marshal(result)
	if err != nil {
		logger.Logger.Error("serialize request info failed! cause: ", err.Error())
	}
	return resBytes, err
}

func NewFileSave(fileName string) *FileSave {
	if isExistFile(fileName) {
		logger.Logger.Fatal(fileName, " is exist, please move old file...")
	}
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Logger.Fatal("cannot open file: ", err.Error())
	}
	return &FileSave{
		file: f,
		iow:  bufio.NewWriter(f),
	}
}

func (fs *FileSave) Save(req *model.Request) {
	reqResultInfo, err := getJsonSerialize(*req)
	if err != nil {
		logger.Logger.Error("cannot serialization")
		return
	}
	_, err = fs.iow.WriteString(string(reqResultInfo) + "\n")
	if err != nil {
		logger.Logger.Error("cannot write to file")
	}
}

func isExistFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}

func (fs *FileSave) Close() {
	fs.iow.Flush()
	fs.file.Close()
}
