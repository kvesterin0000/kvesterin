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
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
			}
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			err := p.tmpl.Lookup("request").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
