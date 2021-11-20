package pages

import (
	"github.com/wasteimage/dist/server/db"
	"github.com/wasteimage/dist/server/pages/locales"
	"net/http"
	"strconv"
	"text/template"
)

type RequestContext struct {
	rw     http.ResponseWriter
	r      *http.Request
	userID int
	theme  string
}

func ContextFromRWR(rw http.ResponseWriter, r *http.Request) RequestContext {
	userID := readSession(r)
	theme := readTheme(r)

	return RequestContext{
		rw:     rw,
		r:      r,
		userID: userID,
		theme:  theme,
	}
}

type PageInfo struct {
	Name   string
	Path   string
	BackTo string
}

type Page interface {
	Info() PageInfo
	Get(rc RequestContext)
	Post(rc RequestContext)
}

type page struct {
	name string

	db   db.Service
	tmpl *template.Template
	loc  *locales.Locales
}

func newPage(
	name string,
	pgService db.Service,
	tmpl *template.Template,
	loc *locales.Locales,
) page {
	return page{
		name: name,
		db:   pgService,
		tmpl: tmpl,
		loc:  loc,
	}
}

func (p *page) Info() PageInfo {
	return PageInfo{
		Name:   p.name,
		Path:   "/" + p.name + "/",
		BackTo: "../" + p.name,
	}
}

func readTheme(r *http.Request) string {
	session, err := r.Cookie(themeCookie)
	if err != nil {
		return ""
	}
	return session.Value
}

func readSession(r *http.Request) int {
	session, err := r.Cookie(sessionCookie)
	if err != nil {
		return 0
	}
	sessionInt, err := strconv.Atoi(session.Value)
	if err != nil {
		return 0
	}
	return sessionInt
}
