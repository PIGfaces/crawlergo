package resultsave

import "crawlergo/pkg/model"

type ResultSave interface {
	Save(*model.Request)
	Close()
}
