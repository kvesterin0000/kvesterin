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
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"),
				"prices_where", "prices_lower", "prices_single", "prices_single_st", "prices_single_p",
				"prices_ep", "prices_ep_st", "prices_ep_p", "prices_album", "prices_album_st", "prices_album_p",
				"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
				"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
			)
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
				"locales":  locales,
			}
			err = p.tmpl.Lookup("prices").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
