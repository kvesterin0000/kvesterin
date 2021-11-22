package pages

import (
	"github.com/wasteimage/dist/server/db"
	"github.com/wasteimage/dist/server/pages/locales"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

type Theme uint

const (
	ThemeWhite Theme = iota
	ThemeBlack
)

type RequestContext struct {
	rw        http.ResponseWriter
	r         *http.Request
	userID    int
	theme     Theme
	themeOpts ThemeOpts
}

func ContextFromRWR(rw http.ResponseWriter, r *http.Request) RequestContext {
	userID := readSession(r)
	theme := readTheme(r)

	rc := RequestContext{
		rw:     rw,
		r:      r,
		userID: userID,
		theme:  theme,
	}
	// FIXME: put all path-including variables to common config structure
	if rc.IsDark() {
		rc.themeOpts = ThemeOpts{
			"style_black.css",
			"logo_white.png",
			"cover_black.png",
			"strelka_white.png",
			"checked",
			ColorSuccess,
		}
	} else {
		rc.themeOpts = ThemeOpts{
			"style.css",
			"logo.png",
			"cover.png",
			"strelka.png",
			"",
			ColorPrimary,
		}
	}
	return rc
}

func (rc *RequestContext) IsDark() bool {
	return rc.theme == ThemeBlack
}

func (rc *RequestContext) IsLoggedIn() bool {
	return rc.userID > 0
}

func (rc *RequestContext) Redirect(name string) {
	http.Redirect(rc.rw, rc.r, GetPage(name).Info().BackTo, http.StatusFound)
}

func (rc *RequestContext) SetCookie(name, value string, duration time.Duration) {
	cookie := http.Cookie{
		Name:    name,
		Value:   value,
		Path:    "/",
		Domain:  "",
		Expires: time.Now().Add(duration),
	}
	http.SetCookie(rc.rw, &cookie)
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

func readTheme(r *http.Request) Theme {
	themeStr, err := r.Cookie(themeCookie)
	if err != nil {
		return 0
	}
	theme, _ := strconv.Atoi(themeStr.Value)
	return Theme(theme)
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
