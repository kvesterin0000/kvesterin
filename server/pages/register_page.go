package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const registerPage = "register"

func init() {
	// Register Page
	initPages = append(initPages, func(p *Pages) Page {
		var pg page
		pg.name = registerPage
		pg.get = func(rw http.ResponseWriter, r *http.Request) {
			userId := readSession(r)
			locales, err := p.loc.TranslatePage(r.Header.Get("Accept-Language"),
				"reg_p", "reg_email", "login_user", "login_pass", "reg_complete", "reg_login",
				"nav_main", "nav_prices", "nav_profile", "nav_cabinet", "nav_request", "nav_logout", "nav_login",
				"footer_info", "footer_vk", "footer_yt", "footer_dev", "footer_more", "footer_dist",
			)
			var params = map[string]interface{}{
				"loggedIn": userId > 0,
				"pages":    p.GetPagesInfo(),
				"locales":  locales,
			}
			if userId > 0 {
				http.Redirect(rw, r, "../cabinet", http.StatusFound)
			}
			err = p.tmpl.Lookup("register").Execute(rw, params)
			if err != nil {
				fmt.Println(err)
			}
		}
		pg.post = func(rw http.ResponseWriter, r *http.Request) {
			email := r.FormValue("email")
			username := r.FormValue("text")
			password := r.FormValue("password")
			err := p.pgService.NewUser(email, username, password)
			userId, err := p.pgService.GetUserId(username, password)
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
		return &pg
	})
}
