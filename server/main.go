package main

import (
	"github.com/wasteimage/dist/server/db"
	"github.com/wasteimage/dist/server/pages"
	"net/http"

	_ "github.com/lib/pq"
)

//TODO: Запретить получение чужого релиза пользователем по переходу в строке

func main() {

	pgDB, err := db.New("user=postgres host=0.0.0.0 port=8889 dbname=fuse sslmode=disable")
	if err != nil {
		panic(err)
	}
	pageHandler := pages.New(pgDB).GetHandler()

	server := http.Server{
		Addr:    ":5565",
		Handler: pageHandler,
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
