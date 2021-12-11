package pages

import (
	"fmt"
	"strconv"
	"time"
)

const loginPageName = "login"

const (
	styleHidden = "display: none;"
	styleShown  = "display: block;"
)

var _ Page = &loginPage{}

type loginPage struct {
	passStyle  string
	loginStyle string
	page
}

func (p *loginPage) Get(rc RequestContext) {
	pgLocs := []string{
		"login_p", "login_user", "login_pass", "login_complete", "login_reg",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	rc.SetCookie(sessionCookie, "", 0)
	locales, err := p.loc.TranslatePage(rc.Language(), pgLocs...)
	var params = map[string]interface{}{
		"loggedIn":   rc.userID > 0,
		"pages":      AllPagesInfo(),
		"locales":    locales,
		"wrong":      p.passStyle,
		"loginStyle": p.loginStyle,
		"themeOpts":  rc.themeOpts,
	}
	p.passStyle = styleHidden
	p.loginStyle = styleHidden
	if rc.IsLoggedIn() {
		rc.Redirect(cabinetPageName)
	}
	err = p.tmpl.Lookup("login").Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *loginPage) Post(rc RequestContext) {
	username := rc.r.FormValue("text")
	password := rc.r.FormValue("password")
	if username == "" || password == "" {
		p.loginStyle = styleShown
		rc.Redirect(loginPageName)
	}
	userId, err := p.db.GetUserId(username, password)
	if err != nil || userId <= 0 && len(username) > 0 && len(password) > 0 {
		p.passStyle = styleShown
		rc.Redirect(loginPageName)
	}
	rc.SetCookie(sessionCookie, strconv.Itoa(userId), time.Hour*48)
	rc.Redirect(cabinetPageName)
	return
}
