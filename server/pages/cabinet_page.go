package pages

import (
	"fmt"
	"net/http"
)

const cabinetPage = "cabinet"

func init() {
	// Cabinet page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = cabinetPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"),
				"cabinet_p", "cabinet_more", "cabinet_upload", "cabinet_no_releases", "status_success",
				"status_pending", "status_default", "status_canceled", "nav_main", "nav_prices", "nav_profile",
				"nav_cabinet", "nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt",
				"footer_dev", "footer_more", "footer_dist",
			)
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
				"pages":    p.GetPagesInfo(),
				"locales":  locales,
			}
			err = p.tmpl.Lookup(pg.name).Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
