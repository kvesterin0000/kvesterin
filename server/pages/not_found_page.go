package pages

import (
	"fmt"
	"net/http"
)

const notFoundPage = "notFound"

func init() {
	// notFound Page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = notFoundPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			err := p.tmpl.Lookup("notFound").Execute(rw, nil)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
