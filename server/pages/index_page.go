package pages

import (
	"fmt"
)

const indexPageName = "index"

var _ Page = &indexPage{}

type indexPage struct {
	page
}

func (p *indexPage) Get(rq RequestContext) {
	var currentTheme string
	var navLogo string
	var colorTheme string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"), "title", "desc", "start",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	if err != nil {
		GetPage(notFoundPageName).Get(rq)
		return
	}
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	err = p.tmpl.Lookup(indexPageName).Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *indexPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
