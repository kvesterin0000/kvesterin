package pages

import (
	"fmt"
)

const requestPageName = "request"

var _ Page = &requestPage{}

type requestPage struct {
	sent    string
	warning string
	page
}

func (p *requestPage) Get(rc RequestContext) {
	pgLocs := []string{
		"request_p", "request_release_name", "request_text", "request_send", "request_success",
		"request_sent", "request_btn_success", "nav_main", "nav_prices", "nav_profile", "nav_cabinet",
		"nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev",
		"footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.r.Header.Get("Accept-Language"), pgLocs...)
	var params = map[string]interface{}{
		"loggedIn":   rc.userID > 0,
		"pages":      AllPagesInfo(),
		"locales":    locales,
		"loginStyle": p.warning,
		"success":    p.sent,
		"themeOpts":  rc.themeOpts,
	}
	p.sent = "display: none;"
	p.warning = "display: none;"
	if rc.userID <= 0 {
		rc.Redirect(loginPageName)
	}
	err = p.tmpl.Lookup("request").Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *requestPage) Post(rc RequestContext) {
	release := rc.r.FormValue("release")
	request := rc.r.FormValue("request")
	if len(release) < 1 || len(request) < 1 {
		p.warning = "display: block;"
		rc.Redirect(requestPageName)
	} else {
		p.sent = "display: block;"
		rc.Redirect(requestPageName)
	}
	return
}
