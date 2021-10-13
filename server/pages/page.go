package pages

import (
	"net/http"
	"strconv"
)

type PageInfo struct {
	Name   string
	Path   string
	BackTo string
}

type Page interface {
	Info() PageInfo
	Get(rw http.ResponseWriter, r *http.Request)
	Post(rw http.ResponseWriter, r *http.Request)
}

type page struct {
	name string
	get  http.HandlerFunc
	post http.HandlerFunc
}

func (p *page) Info() PageInfo {
	return PageInfo{
		Name:   p.name,
		Path:   "/" + p.name + "/",
		BackTo: "../" + p.name,
	}
}

func (p *page) Get(rw http.ResponseWriter, r *http.Request) {
	if p.get != nil {
		p.get(rw, r)
	}
}
func (p *page) Post(rw http.ResponseWriter, r *http.Request) {
	if p.post != nil {
		p.post(rw, r)
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
