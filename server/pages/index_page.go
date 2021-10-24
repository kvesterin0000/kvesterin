package pages

import (
	"fmt"
	"net/http"
)

const indexPage = "index"

func init() {
	// Index page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = indexPage
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
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"), "title", "desc", "start",
				"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
				"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
			)
			if err != nil {
				p.GetPage(notFoundPage).Get(rw, r)
				return
			}
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
				"locales":  locales,
				"theme":    currentTheme,
				"nav_logo": navLogo,
				"color":    colorTheme,
			}
			err = p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
