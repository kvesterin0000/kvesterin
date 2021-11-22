package pages

import (
	"fmt"
	"strings"
)

const uploadPageName = "upload"

var _ Page = &uploadPage{}

type uploadPage struct {
	page
}

func (p *uploadPage) Get(rc RequestContext) {
	var params = map[string]interface{}{
		"loggedIn": rc.userID > 0,
	}

	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	err := p.tmpl.Lookup(uploadPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *uploadPage) Post(rc RequestContext) {
	err := rc.r.ParseForm()
	if err != nil {
		rc.Redirect(cabinetPageName)
	}
	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	pgLocs := []string{
		"cabinet_p", "cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email",
		"settings_btn_submit", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request",
		"nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more",
		"footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.r.Header.Get("Accept-Language"), pgLocs...)
	if err != nil {
		GetPage(notFoundPageName).Get(rc)
		return
	}
	var perfs string
	perfs = strings.Join(rc.r.Form["perf"], ", ")
	releaseName := rc.r.FormValue("releaseName")
	//cover := rc.r.FormValue("cover")
	err = p.db.NewRelease(rc.userID, rc.themeOpts.Cover, releaseName, perfs, "В исполнении")
	if err != nil {
		fmt.Println(err)
	}
	var params = map[string]interface{}{
		"loggedIn":    rc.userID > 0,
		"pages":       AllPagesInfo(),
		"releaseName": releaseName,
		"perfs":       perfs,
		"locales":     locales,
		"themeOpts":   rc.themeOpts,
	}
	err = p.tmpl.Lookup(uploadPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}
