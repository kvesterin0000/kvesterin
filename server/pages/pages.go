package pages

import (
	"fmt"
	"github.com/wasteimage/dist/server/db"
	"golang.org/x/net/html/atom"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const sessionCookie = "sessionId"

type ListPage struct {
	Name string
	Path string
}

type Pages struct {
	pgService db.Service
	tmpl      *template.Template
	list      map[string]Page
}

type Page interface {
	Name() string
	Path() string
	Get(rw http.ResponseWriter, r *http.Request)
	Post(rw http.ResponseWriter, r *http.Request)
}

type page struct {
	pages *Pages
	name  string
	get   http.HandlerFunc
	post  http.HandlerFunc
}

func NewPage(name string) *page {
	return &page{
		name: name,
	}
}

func (p *page) Name() string {
	return p.name
}

func (p *page) Path() string {
	return "/" + p.name + "/"
}

func (p *page) Get(rw http.ResponseWriter, r *http.Request) {
	p.get(rw, r)
}
func (p *page) Post(rw http.ResponseWriter, r *http.Request) {
	p.post(rw, r)
}

func New(service db.Service) *Pages {
	var tmpl = template.Must(template.ParseGlob("pages/*"))
	return &Pages{pgService: service, tmpl: tmpl}
}

func (p *Pages) GetList() []ListPage {
	var listPages []ListPage
	for _, page := range p.list {
		listPages = append(listPages, ListPage{
			Name: page.Name(),
			Path: page.Path(),
		})
	}
	return listPages
}

func (p *Pages) AddPage(page Page) {
	p.list[page.Name()] = page
}

func (p *Pages) GetPage(name string) Page {
	return p.list[name]
}

func (p *Pages) GetHandler() http.Handler {

	mux := http.NewServeMux()
	mux.Handle("/", p)
	//mux.HandleFunc("/", p.indexPage)
	mux.HandleFunc("/cabinet/", p.cabinetPage)
	mux.HandleFunc("/release/", p.releasePage)
	mux.HandleFunc("/login/", p.loginPage)
	mux.HandleFunc("/register/", p.registerPage)
	mux.HandleFunc("/prices/", p.pricesPage)
	//mux.HandleFunc("/request/", p.requestPage)
	//mux.HandleFunc("/upload/", p.uploadPage)
	//mux.HandleFunc("/notFound/", p.notFoundPage)

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
	if r.Method == http.MethodGet {
		page.Get(rw, r)
	} else {
		page.Post(rw, r)
	}
}

func (p *Pages) cabinetPage(rw http.ResponseWriter, r *http.Request) {
	userId := readSession(r)
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
	}

	releases, err := p.pgService.GetReleaseByUserId(userId)
	if err != nil {
		fmt.Println("no releases")
	}
	var params = map[string]interface{}{
		"loggedIn": userId > 0,
		"releases": releases,
	}
	err = tmpl.Lookup("cabinet").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}

}

func (p *Pages) uploadPage(rw http.ResponseWriter, r *http.Request) {
	var perfs []string
	var releaseName string
	if r.Method == http.MethodPost {
		perfs = append(perfs, r.FormValue("perf"))
		releaseName = r.FormValue("releaseName")
	}
	userId := readSession(r)
	var params = map[string]interface{}{
		"loggedIn":    userId > 0,
		"releaseName": releaseName,
		"perfs":       perfs,
	}
	if userId <= 0 {
		http.Redirect(rw, r, "../login", http.StatusSeeOther)
	}
	err := tmpl.Lookup("upload").Execute(rw, params)
	if err != nil {
		fmt.Println(err)
	}
}

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
