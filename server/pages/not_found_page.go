package pages

import (
	"fmt"
)

const notFoundPageName = "notFound"

var _ Page = &notFoundPage{}

type notFoundPage struct {
	page
}

func (p *notFoundPage) Get(rq RequestContext) {
	err := p.tmpl.Lookup(notFoundPageName).Execute(rq.rw, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *notFoundPage) Post(rq RequestContext) {
	p.Get(rq)
}
