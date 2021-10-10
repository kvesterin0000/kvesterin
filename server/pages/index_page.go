package pages

import (
	"fmt"
	"net/http"
)

func (p *Pages) initIndexPage() {
	var name = "index"
	indexPage := NewPage(name)
	indexPage.get = func(rw http.ResponseWriter, r *http.Request) {
		userId := readSession(r)
		var params = map[string]interface{}{
			"loggedIn": userId > 0,
			"pages":    p.GetList(),
		}
		err := p.tmpl.Lookup(name).Execute(rw, params)
		if err != nil {
			fmt.Println(err)
		}
	}
}
