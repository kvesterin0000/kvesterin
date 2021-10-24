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
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			var perfs string
			perfs = strings.Join(r.Form["perf"], ", ")
			releaseName := r.FormValue("releaseName")
			var params = map[string]interface{}{
				"loggedIn":    userId > 0,
				"releaseName": releaseName,
				"perfs":       perfs,
				"theme":       currentTheme,
				"nav_logo":    navLogo,
				"color":       colorTheme,
			}
			err = p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
