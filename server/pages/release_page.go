package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const releasePage = "release"

func init() {
	// Release Page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = releasePage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			releaseIdStr := strings.TrimPrefix(r.RequestURI, "/release/")
			releaseId, err := strconv.Atoi(releaseIdStr)
			if err != nil {
				http.Redirect(rw, r, "../notFound", http.StatusSeeOther)
			}
			userId := readSession(r)
			if userId <= 0 {
				http.Redirect(rw, r, "../login", http.StatusSeeOther)
			}
			tracks, err := p.pgService.GetTrackByReleaseId(releaseId)
			if err != nil {
				fmt.Println("no tracks")
			}
			release, err := p.pgService.GetReleaseById(releaseId)
			if err != nil {
				fmt.Println("can't get release by id")
			}
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"release":  release,
				"tracks":   tracks,
				"pages":    p.GetPagesInfo(),
			}
			err = p.tmpl.Lookup("release").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		return &pg
	})
}
