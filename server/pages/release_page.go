package pages

import (
	"fmt"
	"github.com/wasteimage/dist/server/db"
	"strconv"
	"strings"
)

const releasePageName = "release"

var _ Page = &releasePage{}

type releasePage struct {
	page
}

func (p *releasePage) Get(rc RequestContext) {
	releaseIdStr := strings.TrimPrefix(rc.r.RequestURI, "/release/")
	releaseId, err := strconv.Atoi(releaseIdStr)
	if err != nil {
		rc.Redirect(notFoundPageName)
	}
	pgLocs := []string{
		"release_p", "release_return", "status_success", "status_pending", "status_default",
		"status_canceled", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout",
		"nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.Language(), pgLocs...)
	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	tracks, err := p.db.GetTrackByReleaseId(releaseId)
	if err != nil {
		fmt.Println("no tracks")
	}
	release, err := p.db.GetReleaseById(releaseId)
	if err != nil {
		fmt.Println("can't get release by id")
	}
	releases, err := p.db.GetReleaseByUserId(rc.userID)
	if !findReleaseInReleases(release, releases) {
		rc.Redirect(cabinetPageName)
	}
	var params = map[string]interface{}{
		"loggedIn":  rc.userID > 0,
		"release":   release,
		"tracks":    tracks,
		"pages":     AllPagesInfo(),
		"locales":   locales,
		"themeOpts": rc.themeOpts,
	}
	err = p.tmpl.Lookup("release").Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *releasePage) Post(rq RequestContext) {
	GetPage(notFoundPageName).Get(rq)
}

func findReleaseInReleases(release *db.Release, releases []*db.Release) bool {
	for i := 0; i < len(releases); i++ {
		if releases[i].Name == release.Name {
			return true
		}
	}
	return false
}
