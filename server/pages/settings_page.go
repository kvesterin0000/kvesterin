package pages

import (
	"fmt"
	"strconv"
	"time"
)

const settingsPageName = "settings"

var _ Page = &settingsPage{}

type settingsPage struct {
	page
}

func (p *settingsPage) Get(rc RequestContext) {
	if !rc.IsLoggedIn() {
		rc.Redirect(loginPageName)
	}
	user, err := p.db.GetUser(rc.userID)
	if err != nil {
		fmt.Println(err)
	}
	email := user.Email
	pgLocs := []string{
		"cabinet_p", "cabinet_settings", "settings_change_pass", "settings_old_pass", "settings_new_pass",
		"settings_new_pass2", "settings_btn_change", "settings_email_conf", "settings_email", "settings_web",
		"settings_btn_submit", "settings_btn_save", "nav_main", "nav_prices", "nav_profile", "nav_cabinet",
		"nav_request", "nav_logout", "nav_login", "footer_info", "footer_vk", "footer_yt", "footer_dev",
		"footer_more", "footer_dist",
	}
	locales, err := p.loc.TranslatePage(rc.r.Header.Get("Accept-Language"), pgLocs...)
	if err != nil {
		GetPage(notFoundPageName).Get(rc)
		return
	}
	var params = map[string]interface{}{
		"loggedIn":  rc.userID > 0,
		"pages":     AllPagesInfo(),
		"locales":   locales,
		"email":     email,
		"themeOpts": rc.themeOpts,
	}
	err = p.tmpl.Lookup(settingsPageName).Execute(rc.rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *settingsPage) Post(rc RequestContext) {
	switch rc.r.FormValue("form_name") {
	case "change_password":
		oldPass := rc.r.FormValue("password")
		newPass1 := rc.r.FormValue("password_new1")
		newPass2 := rc.r.FormValue("password_new2")
		if newPass1 == newPass2 {
			err := p.db.UpdatePassword(rc.userID, oldPass, newPass2)
			if err != nil {
				fmt.Println(err)
			}
		}
	case "email_confirmation":
	case "change_theme":
		if rc.r.FormValue("theme") == "on" {
			rc.SetCookie(themeCookie, strconv.Itoa(1), time.Hour*24*30)
		} else {
			rc.SetCookie(themeCookie, strconv.Itoa(0), 0)
		}
	}
	rc.Redirect(settingsPageName)
}
