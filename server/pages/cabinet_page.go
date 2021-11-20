package pages

import (
	"fmt"
	"net/http"
)

const cabinetPageName = "cabinet"

var _ Page = &cabinetPage{}

type cabinetPage struct {
	page
}

func (p *cabinetPage) Get(rq RequestContext) {
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
	locs, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"),
		"cabinet_p", "cabinet_settings", "cabinet_more", "cabinet_upload", "cabinet_no_releases",
		"status_success", "status_pending", "status_default", "status_canceled", "nav_main", "nav_prices",
		"nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk",
		"footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}

	releases, err := p.db.GetReleaseByUserId(rq.userID)
	if err != nil {
		fmt.Println("no releases")
	}
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"releases": releases,
		"pages":    AllPagesInfo(),
		"locales":  locs,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	err = p.tmpl.Lookup(cabinetPageName).Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *cabinetPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
