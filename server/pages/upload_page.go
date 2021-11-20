package pages

import (
	"fmt"
	"net/http"
	"strings"
)

const uploadPageName = "upload"

var _ Page = &uploadPage{}

type uploadPage struct {
	page
}

func (p *uploadPage) Get(rq RequestContext) {
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
	}

	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}
	err := p.tmpl.Lookup(uploadPageName).Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *uploadPage) Post(rq RequestContext) {
	err := rq.r.ParseForm()
	if err != nil {
		http.Redirect(rq.rw, rq.r, "../cabinet", http.StatusSeeOther)
	}
	var currentTheme string
	var navLogo string
	var colorTheme string
	var cover string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
		cover = "cover_black.png"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
		cover = "cover.png"
	}
	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"), "cabinet_p",
		"cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email",
		"settings_btn_submit", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request",
		"nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more",
		"footer_dist",
	)
	if err != nil {
		GetPage(notFoundPageName).Get(rq)
		return
	}
	var perfs string
	perfs = strings.Join(rq.r.Form["perf"], ", ")
	releaseName := rq.r.FormValue("releaseName")
	var params = map[string]interface{}{
		"loggedIn":    rq.userID > 0,
		"pages":       AllPagesInfo(),
		"releaseName": releaseName,
		"perfs":       perfs,
		"theme":       currentTheme,
		"nav_logo":    navLogo,
		"color":       colorTheme,
		"cover":       cover,
		"locales":     locales,
	}
	err = p.tmpl.Lookup(uploadPageName).Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}
