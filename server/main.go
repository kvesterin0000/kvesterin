package main

import (
	"github.com/wasteimage/dist/server/db"
	"github.com/wasteimage/dist/server/pages"
	"github.com/wasteimage/dist/server/pages/locales"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

func main() {

	pgDB, err := db.New("user=postgres host=0.0.0.0 port=8889 dbname=fuse sslmode=disable")
	if err != nil {
		panic(err)
	}

	var tmpl = template.Must(template.ParseGlob("pages/*"))
	loc, err := locales.New("languages/en.json", "languages/ru.json")
	if err != nil {
		panic(err)
	}

	pagesHandler := pages.New(pgDB, tmpl, loc).GetHandler()

	server := http.Server{
		Addr:    ":5565",
		Handler: pagesHandler,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
