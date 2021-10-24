package pages

import (
	"fmt"
	"net/http"
	"strings"
)

func init() {
	// Upload page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = "upload"
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
			}

			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			err := p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		pg.post = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			err := r.ParseForm()
			if err != nil {
				http.Redirect(rw, r, "../cabinet", http.StatusSeeOther)
			}
			var currentTheme string
			var navLogo string
			var colorTheme string
			var cover string
			theme := readTheme(r)
			if theme == "SGreen" {
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
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"), "cabinet_p",
				"cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
				"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email",
				"settings_btn_submit", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request",
				"nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more",
				"footer_dist",
			)
			if err != nil {
				p.GetPage(notFoundPage).Get(rw, r)
				return
			}
			var perfs string
			perfs = strings.Join(r.Form["perf"], ", ")
			releaseName := r.FormValue("releaseName")
			var params = map[string]interface{}{
				"loggedIn":    userId > 0,
				"pages":       p.GetPagesInfo(),
				"releaseName": releaseName,
				"perfs":       perfs,
				"theme":       currentTheme,
				"nav_logo":    navLogo,
				"color":       colorTheme,
				"cover":       cover,
				"locales":     locales,
			}
			err = p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
