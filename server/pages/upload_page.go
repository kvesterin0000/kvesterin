package pages

import (
	"fmt"
	"net/http"
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
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			var perfs string
			perfs = r.FormValue("perf") + ","
			releaseName := r.FormValue("releaseName")
			var params = map[string]interface{}{
				"loggedIn":    userId > 0,
				"releaseName": releaseName,
				"perfs":       perfs,
			}
			err := p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
