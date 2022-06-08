package resultsave

import (
	"strings"

	"github.com/PIGfaces/crawlergo/pkg/logger"
	"github.com/PIGfaces/crawlergo/pkg/model"
)

type SubDomainSave struct {
	*AllDomainSave
	domainLimited string
}

func NewDomainSave(domainFile string, domainLimited string) *SubDomainSave {
	return &SubDomainSave{
		NewAllDomainSave(domainFile),
		domainLimited,
	}
}

func (ds *SubDomainSave) Save(req *model.Request) {
	if !ds.filter(req) {
		if _, err := ds.iow.WriteString(req.URL.Hostname()); err != nil {
			logger.Logger.Error("save all domain failed! ", err.Error())
		}
	}
}

func (ds *SubDomainSave) filter(req *model.Request) bool {
	domain := req.URL.Hostname()
	// 若不是主域名后缀或者已保存就过滤
	if !strings.HasSuffix(domain, ds.domainLimited) || ds.domainSet.Contains(domain) {
		return true
	}
	ds.domainSet.Add(domain)
	return false
}
