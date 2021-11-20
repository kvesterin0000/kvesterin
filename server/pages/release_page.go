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

func (p *releasePage) Get(rw http.ResponseWriter, r *http.Request) {
	releaseIdStr := strings.TrimPrefix(r.RequestURI, "/release/")
	releaseId, err := strconv.Atoi(releaseIdStr)
	if err != nil {
		http.Redirect(rw, r, "../notFound", http.StatusSeeOther)
	}
	userId := readSession(r)
	var currentTheme string
	var navLogo string
	var colorTheme string
	theme := readTheme(r)
	if theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
	}
	locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"),
		"release_p", "release_return", "status_success", "status_pending", "status_default",
		"status_canceled", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout",
		"nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
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
		"loggedIn": userId > 0,
		"release":  release,
		"tracks":   tracks,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	err = p.tmpl.Lookup("release").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *releasePage) Post(rw http.ResponseWriter, r *http.Request) {
	GetPage(notFoundPageName).Get(rw, r)
}
