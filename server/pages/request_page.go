package pages

import (
	"fmt"
	"net/http"
)

const requestPage = "request"

func init() {
	// request Page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		var sent string
		var warning string
		pg.name = requestPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			var currentTheme string
			var navLogo string
			var colorTheme string
			theme := readTheme(r)
			if theme == "SGreen" {
				currentTheme = "style_black.css"
				navLogo = "logo_white.png"
				colorTheme = "success"
			} else {
				currentTheme = "style.css"
				navLogo = "logo.png"
				colorTheme = "primary"
			}
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"),
				"request_p", "request_release_name", "request_text", "request_send", "request_success",
				"request_sent", "request_btn_success", "nav_main", "nav_prices", "nav_profile", "nav_cabinet",
				"nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev",
				"footer_more", "footer_dist",
			)
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
				"locales":  locales,
				"theme":    currentTheme,
				"nav_logo": navLogo,
				"color":    colorTheme,
				"warning":  warning,
				"success":  sent,
			}
			sent = "display: none;"
			warning = "display: none;"
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			err = p.tmpl.Lookup("request").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		pg.post = func(rw http.ResponseWriter, r *http.Request) {
			release := r.FormValue("release")
			request := r.FormValue("request")
			if len(release) < 1 || len(request) < 1 {
				warning = "display: block;"
				http.Redirect(rw, r, "/request/", http.StatusFound)
			} else {
				sent = "display: block;"
				http.Redirect(rw, r, "/request/", http.StatusFound)
			}
			return
		}
		return &pg
	})
}
