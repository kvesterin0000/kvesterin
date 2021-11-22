package pages

import (
	"fmt"
	"net/http"
)

const cabinetPageName = "cabinet"

type Color string

const (
	ColorSuccess Color = "success"
	ColorPrimary Color = "primary"
)

type ThemeOpts struct {
	CSS           string
	Logo          string
	Cover         string
	Pointer       string
	CheckboxValue string
	BtnColor      Color
}

var _ Page = &cabinetPage{}

type cabinetPage struct {
	page
}

func (p *cabinetPage) Get(rc RequestContext) {
	pgLocs := []string{
		"cabinet_p", "cabinet_settings", "cabinet_more", "cabinet_upload", "cabinet_no_releases",
		"status_success", "status_pending", "status_default", "status_canceled", "nav_main", "nav_prices",
		"nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk",
		"footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}

	releases, err := p.db.GetReleaseByUserId(rc.userID)
	if err != nil {
		fmt.Println("no releases")
	}
	locales, err := p.loc.TranslatePage(rc.r.Header.Get("Accept-Language"), pgLocs...)
	var params = map[string]interface{}{
		"loggedIn":  rc.IsLoggedIn(),
		"releases":  releases,
		"pages":     AllPagesInfo(),
		"locales":   locales,
		"themeOpts": rc.themeOpts,
	}
	err = p.tmpl.Lookup(cabinetPageName).Execute(rc.rw, params)
	if err != nil {
		nf := GetPage(notFoundPageName)
		if nf == nil {
			rc.rw.WriteHeader(http.StatusInternalServerError)
		} else {
			nf.Get(rc)
		}
	}
}

func (p *cabinetPage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
