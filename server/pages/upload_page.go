package pages

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
)

const uploadPageName = "upload"
const coverUploadPath = "resources/release covers/"
const formCover = "cover"

var _ Page = &uploadPage{}

type uploadPage struct {
	page
}

func (p *uploadPage) Get(rc RequestContext) {
	var params = map[string]interface{}{
		"loggedIn": rc.userID > 0,
	}

	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	err := p.tmpl.Lookup(uploadPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *uploadPage) Post(rc RequestContext) {
	err := rc.r.ParseForm()
	if err != nil {
		rc.Redirect(cabinetPageName)
	}
	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	pgLocs := []string{
		"cabinet_p", "cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email",
		"settings_btn_submit", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request",
		"nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more",
		"footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.Language(), pgLocs...)
	if err != nil {
		GetPage(notFoundPageName).Get(rc)
		return
	}
	err = rc.r.ParseMultipartForm(1 << 20)
	if err != nil {
		fmt.Println(err)
		return
	}
	var perfs string
	perfs = strings.Join(rc.r.Form["perf"], ", ")
	releaseName := rc.r.FormValue("releaseName")
	fileName, err := parseCover(rc.r)
	if err != nil {
		rc.rw.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	releases, err := p.db.GetReleaseByUserId(rc.userID)
	if !FindReleaseInReleases(releaseName, perfs, releases) {
		err = p.db.NewRelease(rc.userID, fileName, releaseName, perfs, "В исполнении")
		if err != nil {
			fmt.Println(err)
		}
	}
	releaseId, err := p.db.GetReleaseId(releaseName, perfs)
	release, err := p.db.GetReleaseById(releaseId)
	if err != nil {
		fmt.Println("can't get release by id")
	}
	tracks, err := p.db.GetTrackByReleaseId(releaseId)
	if err != nil {
		fmt.Println("can't get tracks by release id")
	}
	var params = map[string]interface{}{
		"loggedIn":    rc.userID > 0,
		"pages":       AllPagesInfo(),
		"releaseName": releaseName,
		"perfs":       perfs,
		"locales":     locales,
		"themeOpts":   rc.themeOpts,
		"release":     release,
		"tracks":      tracks,
	}
	err = p.tmpl.Lookup(uploadPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func parseCover(r *http.Request) (string, error) {
	cover, fileHeader, err := r.FormFile(formCover)
	if err != nil {
		return "", err
	}
	fileName := path.Join(coverUploadPath, fileHeader.Filename)
	f, err := os.Create(fileName)
	if err != nil {
		return "", err
	}
	buf := bufio.NewReader(cover)
	_, err = buf.WriteTo(f)
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}
	err = cover.Close()
	if err != nil {
		return "", err
	}
	return fileHeader.Filename, nil
}
