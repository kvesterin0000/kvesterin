package pages

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (p *Pages) releasePage(rw http.ResponseWriter, r *http.Request) {
	releaseIdStr := strings.TrimPrefix(r.RequestURI, "/release/")
	releaseId, err := strconv.Atoi(releaseIdStr)
	if err != nil {
		http.Redirect(rw, r, "../notFound", http.StatusSeeOther)
	}
	userId := readSession(r)
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
	}
	tracks, err := p.pgService.GetTrackByReleaseId(releaseId)
	if err != nil {
		fmt.Println("no tracks")
	}
	release, err := p.pgService.GetReleaseById(releaseId)
	if err != nil {
		fmt.Println("can't get release by id")
	}
	var params = map[string]interface{}{
		"loggedIn": userId > 0,
		"release":  release,
		"tracks":   tracks,
	}
	err = tmpl.Lookup("release").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}

}

func (p *Pages) loginPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		trashCookie := http.Cookie{
			Name:    sessionCookie,
			Path:    "/",
			Expires: time.Now(),
		}
		http.SetCookie(rw, &trashCookie)
	}
	if r.Method == http.MethodPost {
		username := r.FormValue("text")
		password := r.FormValue("password")
		userId, err := p.pgService.GetUserId(username, password)
		if err != nil || userId <= 0 {
			http.Redirect(rw, r, "/login/", http.StatusFound)
		}
		session := http.Cookie{
			Name:    sessionCookie,
			Value:   strconv.Itoa(userId),
			Path:    "/",
			Domain:  "",
			Expires: time.Now().Add(time.Hour * 48),
		}
		r.AddCookie(&session)
		http.SetCookie(rw, &session)
		http.Redirect(rw, r, "../cabinet/", http.StatusFound)
		return
	}
	userId := readSession(r)
	var params = map[string]bool{
		"loggedIn": userId > 0,
	}
	if userId > 0 {
		http.Redirect(rw, r, "../cabinet", http.StatusFound)
	}
	err := tmpl.Lookup("login").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}
func (p *Pages) registerPage(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
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
	userId := readSession(r)
	var params = map[string]bool{
		"loggedIn": userId > 0,
	}
	if userId > 0 {
		http.Redirect(rw, r, "../cabinet", http.StatusFound)
	}
	err := tmpl.Lookup("register").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Pages) pricesPage(rw http.ResponseWriter, r *http.Request) {
	userId := readSession(r)
	var params = map[string]bool{
		"loggedIn": userId > 0,
	}
	err := tmpl.Lookup("prices").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Pages) requestPage(rw http.ResponseWriter, r *http.Request) {
	userId := readSession(r)
	var params = map[string]bool{
		"loggedIn": userId > 0,
	}
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
	}
	err := tmpl.Lookup("request").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *Pages) notFoundPage(rw http.ResponseWriter, r *http.Request) {
	err := tmpl.Lookup("notFound").Execute(rw, nil)
	if err != nil {
		fmt.Println(err)
	}
}
