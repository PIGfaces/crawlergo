package resultsave

import "github.com/PIGfaces/crawlergo/pkg/model"

type ResultSave interface {
	Save(*model.Request)
	Close()
}
