package pages

import (
	"fmt"
	"net/http"
)

const pricesPage = "prices"

func init() {
	// Prices Page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = pricesPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
			}
			err := p.tmpl.Lookup("prices").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
