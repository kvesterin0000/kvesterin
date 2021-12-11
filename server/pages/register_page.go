package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const registerPageName = "register"

var _ Page = &registerPage{}

type registerPage struct {
	wrong   string
	warning string
	page
}

func (p *registerPage) Get(rc RequestContext) {
	pgLocs := []string{
		"reg_p", "reg_email", "login_user", "login_pass", "reg_complete", "reg_login",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.Language(), pgLocs...)
	var params = map[string]interface{}{
		"loggedIn":   rc.userID > 0,
		"pages":      AllPagesInfo(),
		"locales":    locales,
		"wrong":      p.wrong,
		"loginStyle": p.warning,
		"themeOpts":  rc.themeOpts,
	}
	p.wrong = "display: none;"
	p.warning = "display: none;"
	if rc.userID > 0 {
		http.Redirect(rc.rw, rc.r, "../cabinet", http.StatusFound)
	}
	err = p.tmpl.Lookup("register").Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *registerPage) Post(rc RequestContext) {
	email := rc.r.FormValue("email")
	username := rc.r.FormValue("text")
	password := rc.r.FormValue("password")
	if email == "" || username == "" || password == "" {
		p.warning = "display: block;"
		rc.Redirect(registerPageName)
		return
	}
	if len(password) < 8 {
		p.wrong = "display: block;"
		rc.Redirect(registerPageName)
		return
	}
	err := p.db.NewUser(email, username, password)
	userId, err := p.db.GetUserId(username, password)
	if err != nil || rc.IsLoggedIn() {
		rc.Redirect(registerPageName)
	}
	session := http.Cookie{
		Name:    sessionCookie,
		Value:   strconv.Itoa(userId),
		Path:    "/",
		Domain:  "*",
		Expires: time.Now().Add(time.Hour * 48),
	}
	http.SetCookie(rc.rw, &session)
	rc.Redirect(cabinetPageName)
	return
}
