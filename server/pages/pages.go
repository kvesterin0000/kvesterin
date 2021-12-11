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
	themeCookie   = "theme"
	langCookie    = "lang"
)

var allPages map[string]Page

func GetPage(name string) Page {
	pg, ok := allPages[name]
	if !ok {
		pg, ok = allPages[notFoundPageName]
		if !ok {
			return nil
		}
	}
	return pg
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
	allPages = make(map[string]Page)
	allPages[cabinetPageName] = &cabinetPage{page: newPage(cabinetPageName, service, tmpl, loc)}
	allPages[indexPageName] = &indexPage{page: newPage(indexPageName, service, tmpl, loc)}
	allPages[loginPageName] = &loginPage{page: newPage(loginPageName, service, tmpl, loc)}
	allPages[notFoundPageName] = &notFoundPage{page: newPage(notFoundPageName, service, tmpl, loc)}
	allPages[pricesPageName] = &pricesPage{page: newPage(pricesPageName, service, tmpl, loc)}
	allPages[registerPageName] = &registerPage{page: newPage(registerPageName, service, tmpl, loc)}
	allPages[releasePageName] = &releasePage{page: newPage(releasePageName, service, tmpl, loc)}
	allPages[requestPageName] = &requestPage{page: newPage(requestPageName, service, tmpl, loc)}
	allPages[settingsPageName] = &settingsPage{page: newPage(settingsPageName, service, tmpl, loc)}
	allPages[uploadPageName] = &uploadPage{page: newPage(uploadPageName, service, tmpl, loc)}
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
	rq := ContextFromRWR(rw, r)

	pageName := parsePageName(r.RequestURI)
	pg := GetPage(pageName)
	if pg == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	switch r.Method {
	case http.MethodGet:
		pg.Get(rq)
	case http.MethodPost:
		pg.Post(rq)
	default:
		GetPage(notFoundPageName).Get(rq)
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
