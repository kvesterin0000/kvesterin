package pages

const newPageName = "new_page"

var _ Page = &newPageType{}

type newPageType struct {
	page
}

func (p *newPageType) Get(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}

func (p *newPageType) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
