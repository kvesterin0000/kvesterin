package pages

import (
	"fmt"
	"net/http"
	"time"
)

const settingsPageName = "settings"

var _ Page = &settingsPage{}

type settingsPage struct {
	page
}

func (p *settingsPage) Get(rw http.ResponseWriter, r *http.Request) {
	userId := readSession(r)
	var currentTheme string
	var navLogo string
	var colorTheme string
	var val string
	theme := readTheme(r)
	if theme == "SGreen" {
		currentTheme = "style_black.css"
		navLogo = "logo_white.png"
		colorTheme = "success"
		val = "checked"
	} else {
		currentTheme = "style.css"
		navLogo = "logo.png"
		colorTheme = "primary"
		val = ""
	}
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
	}
	user, err := p.db.GetUser(userId)
	if err != nil {
		fmt.Println(err)
	}
	email := user.Email
	locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"), "cabinet_p",
		"cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email",
		"settings_btn_submit", "nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request",
		"nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more",
		"footer_dist",
	)
	if err != nil {
		GetPage(notFoundPageName).Get(rw, r)
		return
	}
	var params = map[string]interface{}{
		"loggedIn": userId > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"email":    email,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
		"val":      val,
	}
	err = p.tmpl.Lookup(settingsPageName).Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *settingsPage) Post(rw http.ResponseWriter, r *http.Request) {
	userId := readSession(r)
	switch r.FormValue("form_name") {
	case "change_password":
		oldPass := r.FormValue("password")
		newPass1 := r.FormValue("password_new1")
		newPass2 := r.FormValue("password_new2")
		if newPass1 == newPass2 {
			err := p.db.UpdatePassword(userId, oldPass, newPass2)
			if err != nil {
				fmt.Println(err)
			}
		}
	case "email_confirmation":
	case "change_theme":
		if r.FormValue("theme") == "on" {
			session := http.Cookie{
				Name:    themeCookie,
				Value:   "SGreen",
				Path:    "/",
				Domain:  "",
				Expires: time.Now().Add(time.Hour * 730),
			}
			r.AddCookie(&session)
			http.SetCookie(rw, &session)
		} else {
			trashCookie := http.Cookie{
				Name:    themeCookie,
				Path:    "/",
				Expires: time.Now(),
			}
			http.SetCookie(rw, &trashCookie)
		}
	}
	fmt.Println(r.FormValue("theme"))
	http.Redirect(rw, r, "../settings/", http.StatusFound)
}
