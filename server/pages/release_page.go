package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const releasePageName = "release"

var _ Page = &releasePage{}

type releasePage struct {
	page
}

func (p *releasePage) Get(rq RequestContext) {
	releaseIdStr := strings.TrimPrefix(rq.r.RequestURI, "/release/")
	releaseId, err := strconv.Atoi(releaseIdStr)
	if err != nil {
		http.Redirect(rq.rw, rq.r, "../notFound", http.StatusSeeOther)
	}
	var currentTheme string
	var navLogo string
	var colorTheme string
	if rq.theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
	}
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"),
		"release_p", "release_return", "status_success", "status_pending", "status_default",
		"status_canceled", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout",
		"nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}
	tracks, err := p.db.GetTrackByReleaseId(releaseId)
	if err != nil {
		fmt.Println("no tracks")
	}
	release, err := p.db.GetReleaseById(releaseId)
	if err != nil {
		fmt.Println("can't get release by id")
	}
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"release":  release,
		"tracks":   tracks,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	err = p.tmpl.Lookup("release").Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *releasePage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}
