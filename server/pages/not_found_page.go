package pages

import (
	"fmt"
	"net/http"
)

const notFoundPageName = "notFound"

var _ Page = &notFoundPage{}

type notFoundPage struct {
	page
}

func (p *notFoundPage) Get(rw http.ResponseWriter, r *http.Request) {
	err := p.tmpl.Lookup(notFoundPageName).Execute(rw, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *notFoundPage) Post(rw http.ResponseWriter, r *http.Request) {
	p.Get(rw, r)
}
