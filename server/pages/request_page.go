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
		pg.name = requestPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
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
			}
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			err = p.tmpl.Lookup("request").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
