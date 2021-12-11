package pages

import (
	"fmt"
)

const indexPageName = "index"

var _ Page = &indexPage{}

type indexPage struct {
	page
}

func (p *indexPage) Get(rc RequestContext) {
	pgLocs := []string{
		"title", "desc", "start", "welcome_message1", "welcome_message2", "welcome_message3", "welcome_btn", "nav_main",
		"nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login", "footer_info",
		"footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.Language(), pgLocs...)
	if err != nil {
		GetPage(notFoundPageName).Get(rc)
		return
	}
	var params = map[string]interface{}{
		"loggedIn":  rc.userID > 0,
		"pages":     AllPagesInfo(),
		"locales":   locales,
		"themeOpts": rc.themeOpts,
	}
	err = p.tmpl.Lookup(indexPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *indexPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
