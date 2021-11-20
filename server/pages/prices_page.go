package pages

import (
	"fmt"
)

const pricesPageName = "prices"

var _ Page = &pricesPage{}

type pricesPage struct {
	page
}

func (p *pricesPage) Get(rq RequestContext) {
	var currentTheme string
	var navLogo string
	var colorTheme string
	var pointer string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
		pointer = "strelka_white.png"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
		pointer = "strelka.png"
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"),
		"prices_where", "prices_lower", "prices_single", "prices_single_st", "prices_single_p",
		"prices_ep", "prices_ep_st", "prices_ep_p", "prices_album", "prices_album_st", "prices_album_p",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
		"pointer":  pointer,
	}
	err = p.tmpl.Lookup("prices").Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *pricesPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
