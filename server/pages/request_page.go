package pages

import (
	"fmt"
	"net/http"
)

const requestPageName = "request"

var _ Page = &requestPage{}

type requestPage struct {
	sent    string
	warning string
	page
}

func (p *requestPage) Get(rq RequestContext) {
	var currentTheme string
	var navLogo string
	var colorTheme string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
	} else {
		currentTheme = "style.css"
		// FIXME: theme related stuff must be in config
		navLogo = "logo.png"
		colorTheme = "primary"
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"),
		"request_p", "request_release_name", "request_text", "request_send", "request_success",
		"request_sent", "request_btn_success", "nav_main", "nav_prices", "nav_profile", "nav_cabinet",
		"nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev",
		"footer_more", "footer_dist",
	)
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
		"warning":  p.warning,
		"success":  p.sent,
	}
	p.sent = "display: none;"
	p.warning = "display: none;"
	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}
	err = p.tmpl.Lookup("request").Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *requestPage) Post(rq RequestContext) {
	release := rq.r.FormValue("release")
	request := rq.r.FormValue("request")
	if len(release) < 1 || len(request) < 1 {
		p.warning = "display: block;"
		http.Redirect(rq.rw, rq.r, "/request/", http.StatusFound)
	} else {
		p.sent = "display: block;"
		http.Redirect(rq.rw, rq.r, "/request/", http.StatusFound)
	}
	return
}
