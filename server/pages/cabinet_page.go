package pages

import (
	"fmt"
	"net/http"
)

func init() {
	// Cabinet page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = "cabinet"
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}

			releases, err := p.pgService.GetReleaseByUserId(userId)
			if err != nil {
				fmt.Println("no releases")
			}
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"releases": releases,
			}
			err = p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
