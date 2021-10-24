package pages

import (
	"github.com/wasteimage/dist/server/pages/locales"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/wasteimage/dist/server/db"
)

const (
	sessionCookie = "sessionId"
	themeCookie   = "SpecGreen"
)

var initPages []pageIniter

type pageIniter func(*Pages) Page

type Pages struct {
	pgService db.Service
	tmpl      *template.Template
	loc       *locales.Locales
	list      map[string]Page
}

func New(service db.Service) *Pages {
	var tmpl = template.Must(template.ParseGlob("pages/*"))
	loc, err := locales.New("languages/en.json", "languages/ru.json")
	if err != nil {
		panic(err)
	}
	p := &Pages{pgService: service, tmpl: tmpl, loc: loc}
	for _, initPage := range initPages {
		pg := initPage(p)
		p.AddPage(pg)
	}
	return p
}

func (p *Pages) AddPage(page Page) {
	if p.list == nil {
		p.list = make(map[string]Page)
	}
	p.list[page.Info().Name] = page
}

func (p *Pages) GetPage(name string) Page {
	return p.list[name]
}

func (p *Pages) GetPagesInfo() map[string]PageInfo {
	var pagesInfo = make(map[string]PageInfo)
	for _, page := range p.list {
		pagesInfo[page.Info().Name] = page.Info()
	}
	return pagesInfo
}

func (p *Pages) GetHandler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", p)

	fs := http.FileServer(http.Dir("./resources/"))
	mux.Handle("/resources/", http.StripPrefix("/resources/", fs))
	return mux
}

func (p *Pages) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	var (
		root     = "/"
		pageName string
	)
	if r.RequestURI == root {
		pageName = indexPage
	} else {
		uri := strings.SplitN(r.RequestURI, "/", 3)
		if len(uri) < 2 {
			return
		}
		pageName = uri[1]
	}
	page := p.GetPage(pageName)
	if page == nil {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	switch r.Method {
	case http.MethodGet:
		page.Get(rw, r)
	case http.MethodPost:
		page.Post(rw, r)
	default:
		err := p.tmpl.ExecuteTemplate(rw, "notFound", nil)
		if err != nil {
			log.Fatalf("not found template execution error: %v", err)
		}
	}
}
