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

func (p *settingsPage) Get(rq RequestContext) {
	var currentTheme string
	var navLogo string
	var colorTheme string
	var val string
	if rq.theme == "SGreen" {
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
	if rq.userID <= 0 {
		http.Redirect(rq.rw, rq.r, "../login", http.StatusSeeOther)
	}
	user, err := p.db.GetUser(rq.userID)
	if err != nil {
		fmt.Println(err)
	}
	email := user.Email
	locales, err := p.loc.TranslatePage(rq.r.Header.Get("Accept-Language"), "cabinet_p",
		"cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email", "settings_web",
		"settings_btn_submit", "settings_btn_save", "nav_main", "nav_prices", "nav_profile", "nav_cabinet",
		"nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev",
		"footer_more", "footer_dist",
	)
	if err != nil {
		GetPage(notFoundPageName).Get(rq)
		return
	}
	var params = map[string]interface{}{
		"loggedIn": rq.userID > 0,
		"pages":    AllPagesInfo(),
		"locales":  locales,
		"email":    email,
		"theme":    currentTheme,
		"nav_logo": navLogo,
		"color":    colorTheme,
		"val":      val,
	}
	err = p.tmpl.Lookup(settingsPageName).Execute(rq.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *settingsPage) Post(rq RequestContext) {
	switch rq.r.FormValue("form_name") {
	case "change_password":
		oldPass := rq.r.FormValue("password")
		newPass1 := rq.r.FormValue("password_new1")
		newPass2 := rq.r.FormValue("password_new2")
		if newPass1 == newPass2 {
			err := p.db.UpdatePassword(rq.userID, oldPass, newPass2)
			if err != nil {
				fmt.Println(err)
			}
		}
	case "email_confirmation":
	case "change_theme":
		if rq.r.FormValue("theme") == "on" {
			session := http.Cookie{
				Name:    themeCookie,
				Value:   "SGreen",
				Path:    "/",
				Domain:  "",
				Expires: time.Now().Add(time.Hour * 730),
			}
			rq.r.AddCookie(&session)
			http.SetCookie(rq.rw, &session)
		} else {
			// FIXME: make method for zero time cookie
			trashCookie := http.Cookie{
				Name:    themeCookie,
				Path:    "/",
				Expires: time.Now(),
			}
			http.SetCookie(rq.rw, &trashCookie)
		}
	}
	fmt.Println(rq.r.FormValue("theme"))
	http.Redirect(rq.rw, rq.r, "../settings/", http.StatusFound)
}
