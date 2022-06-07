package resultsave

import (
	"crawlergo/pkg/logger"
	"crawlergo/pkg/model"

	mapset "github.com/deckarep/golang-set"
)

type AllDomainSave struct {
	*FileSave
	domainSet mapset.Set
}

func NewAllDomainSave(domainFile string) *AllDomainSave {
	return &AllDomainSave{
		NewFileSave(domainFile),
		mapset.NewSet(),
	}
}

func (ds *AllDomainSave) Save(req *model.Request) {
	if !ds.filter(req) {
		if _, err := ds.iow.WriteString(req.URL.Hostname()); err != nil {
			logger.Logger.Error("save all domain failed! ", err.Error())
		}

	}
}

func (ds *AllDomainSave) filter(req *model.Request) bool {
	domain := req.URL.Hostname()
	if ds.domainSet.Contains(domain) {
		return true
	}
	ds.domainSet.Add(domain)
	return false
}
