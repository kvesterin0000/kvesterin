package pages

import (
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/wasteimage/dist/server/db"
)

const sessionCookie = "sessionId"

var initPages []pageIniter

type pageIniter func(*Pages) Page

type Pages struct {
	pgService db.Service
	tmpl      *template.Template
	list      map[string]Page
}

func New(service db.Service) *Pages {
	var tmpl = template.Must(template.ParseGlob("pages/*"))
	p := &Pages{pgService: service, tmpl: tmpl}
	for _, initPage := range initPages {
		pg := initPage(p)
		p.AddPage(pg)
	}
	return p
}

func (p *Pages) AddPage(page Page) {
	p.list[page.Info().Name] = page
}

func (p *Pages) GetPage(name string) Page {
	return p.list[name]
}

func (p *Pages) GetPagesInfo() []PageInfo {
	var pagesInfo []PageInfo
	for _, page := range p.list {
		pagesInfo = append(pagesInfo, page.Info())
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
	uri := strings.SplitN(r.RequestURI, "/", 3)
	if len(uri) != 3 {
		return
	}
	pageName := uri[1]
	page := p.GetPage(pageName)
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
