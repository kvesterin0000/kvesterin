package pages

import (
	"net/http"
	"strings"
	"text/template"

	"github.com/wasteimage/dist/server/db"
	"github.com/wasteimage/dist/server/pages/locales"
)

const (
	sessionCookie = "sessionId"
	themeCookie   = "SpecGreen"
)

var allPages map[string]Page

func GetPage(name string) Page {
	return allPages[name]
}

func AllPagesInfo() map[string]PageInfo {
	var pagesInfo = make(map[string]PageInfo)
	for _, page := range allPages {
		pagesInfo[page.Info().Name] = page.Info()
	}
	return pagesInfo
}

type Pages struct{}

func New(service db.Service, tmpl *template.Template, loc *locales.Locales) *Pages {
	allPages[cabinetPageName] = &cabinetPage{newPage(cabinetPageName, service, tmpl, loc)}
	return &Pages{}
}

func (p *Pages) GetHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", p)

	fs := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/", http.StripPrefix("/resources/", fs))
	return mux
}

func (p *Pages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	pageName := parsePageName(r.RequestURI)
	pg := GetPage(pageName)
	if pg == nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	switch r.Method {
	case http.MethodGet:
		pg.Get(rw, r)
	case http.MethodPost:
		pg.Post(rw, r)
	default:
		GetPage(notFoundPageName).Get(rw, r)
	}
}

func parsePageName(uri string) string {
	var (
		root     = "/"
		pageName string
	)
	if uri == root {
		pageName = indexPageName
	} else {
		uri := strings.SplitN(uri, "/", 3)
		if len(uri) < 2 {
			return notFoundPageName
		}
		pageName = uri[1]
	}
	return pageName
}
