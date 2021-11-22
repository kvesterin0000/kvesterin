package pages

import (
	"fmt"
)

const pricesPageName = "prices"

var _ Page = &pricesPage{}

type pricesPage struct {
	page
}

func (p *pricesPage) Get(rc RequestContext) {
	pgLocs := []string{
		"prices_where", "prices_lower", "prices_single", "prices_single_st", "prices_single_p",
		"prices_ep", "prices_ep_st", "prices_ep_p", "prices_album", "prices_album_st", "prices_album_p",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.r.Header.Get("Accept-Language"), pgLocs...)
	var params = map[string]interface{}{
		"loggedIn":  rc.userID > 0,
		"pages":     AllPagesInfo(),
		"locales":   locales,
		"themeOpts": rc.themeOpts,
	}
	err = p.tmpl.Lookup("prices").Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *pricesPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
