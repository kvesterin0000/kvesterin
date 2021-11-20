package pages

import "net/http"

const newPageName = "new_page"

var _ Page = &newPageType{}

type newPageType struct {
	page
}

func (p *newPageType) Get(rw http.ResponseWriter, r *http.Request) {
	GetPage(notFoundPageName).Get(rw, r)
}

func (p *newPageType) Post(rw http.ResponseWriter, r *http.Request) {
	GetPage(notFoundPageName).Get(rw, r)
}
