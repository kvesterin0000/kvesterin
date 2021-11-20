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

func (p *registerPage) Get(rw http.ResponseWriter, r *http.Request) {
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
		"reg_p", "reg_email", "login_user", "login_pass", "reg_complete", "reg_login",
		"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
		"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
	)
	var params = map[string]interface{}{
		"loggedIn": userId > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"wrong":    p.wrong,
		"warning":  p.warning,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
	}
	p.wrong = "display: none;"
	p.warning = "display: none;"
	if userId > 0 {
		http.Redirect(rw, r, "../cabinet", http.StatusFound)
	}
	err = p.tmpl.Lookup("register").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *registerPage) Post(rw http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	username := r.FormValue("text")
	password := r.FormValue("password")
	if len(email) < 1 || len(username) < 1 || len(password) < 1 {
		p.warning = "display: block;"
		http.Redirect(rw, r, "/register/", http.StatusFound)
		return
	}
	if len(password) < 8 {
		p.wrong = "display: block;"
		http.Redirect(rw, r, "/register/", http.StatusFound)
		return
	}
	err := p.db.NewUser(email, username, password)
	userId, err := p.db.GetUserId(username, password)
	if err != nil || userId <= 0 {
		http.Redirect(rw, r, "/register/", http.StatusFound)
	}
	session := http.Cookie{
		Name:    sessionCookie,
		Value:   strconv.Itoa(userId),
		Path:    "/",
		Domain:  "*",
		Expires: time.Now().Add(time.Hour * 48),
	}
	r.AddCookie(&session)
	http.SetCookie(rw, &session)
	http.Redirect(rw, r, "../cabinet", http.StatusFound)
	return
}
