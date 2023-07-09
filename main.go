package main

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
)

var tmpl *template.Template

type BigBoy struct {
	Todos   map[uuid.UUID]todo
	Todones map[uuid.UUID]todo
}

func getIndex(w http.ResponseWriter, _ *http.Request) {
	err := tmpl.Execute(w, BigBoy{
		Todos:   Todos,
		Todones: Todones,
	})
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	tmpl, err = template.ParseFiles("./web/base.gohtml")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err = tmpl.ParseGlob("./web/templates/*.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", getIndex)
	http.HandleFunc("/todos/", handleTodos)

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/robots.txt")
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))
	log.Fatal(http.ListenAndServe(":8000", nil))
}
