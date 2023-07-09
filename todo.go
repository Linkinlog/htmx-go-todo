package main

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type todo struct {
	Id      uuid.UUID
	Name    string
	Details string
}

var Todos = make(map[uuid.UUID]todo)
var Todones = make(map[uuid.UUID]todo)

var todoHTML = template.Must(template.ParseFiles("./web/templates/todo.gohtml", "./web/templates/todone.gohtml", "./web/templates/todosList.gohtml"))

func init() {
	id1 := uuid.New()
	Todos[id1] = todo{
		Id:      id1,
		Name:    "Expand the family",
		Details: "Get more puppies",
	}
	id2 := uuid.New()
	Todos[id2] = todo{
		Id:      id2,
		Name:    "Learn CSS",
		Details: "Center more divs",
	}
	id3 := uuid.New()
	Todos[id3] = todo{
		Id:      id3,
		Name:    "Learn HTML",
		Details: "Create more divs",
	}
}

func newTodo(name, details string) *todo {
	id := uuid.New()
	newTodoItem := &todo{
		Id:      id,
		Name:    name,
		Details: details,
	}
	Todos[id] = *newTodoItem
	return newTodoItem
}

func uuidFromURL(url url.URL) uuid.UUID {
	pathParts := strings.Split(url.Path, "/")

	// Ensure that there is an ID at the expected path part
	if len(pathParts) < 3 {
		log.Fatal("not enough path parts")
	}

	// Extract the ID part
	return uuid.MustParse(pathParts[2])

}

func handleTodos(w http.ResponseWriter, r *http.Request) {
	// If the HX-Request header isn't set to true,
	// we aren't working with HTMX, so just redirect home.
	if r.Header.Get("HX-Request") != "true" {
		http.Redirect(w, r, "/", 302)
	}

	if r.Body != nil {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}
	}

	if r.Method == "POST" && r.FormValue("_method") == "PATCH" {
		updateTodo(w, r)
	} else if r.Method == "POST" && r.FormValue("_method") == "DELETE" {
		delete(Todones, uuidFromURL(*r.URL))
		reloadTodos(w)
	} else if r.Method == "POST" {
		addTodo(w, r)
	}
}

func updateTodo(w http.ResponseWriter, r *http.Request) {
	id := uuidFromURL(*r.URL)
	// Need to check both maps for now, IDK
	_, ok := Todos[id]
	_, found := Todones[id]
	if !ok && !found {
		http.Error(w, "Couldn't delete", http.StatusBadRequest)
		return
	}

	if r.FormValue("complete") == "true" {
		Todones[id] = Todos[id]
		delete(Todos, id)
	}
	reloadTodos(w)
}

func addTodo(w http.ResponseWriter, r *http.Request) {
	name := r.Form.Get("name")
	details := r.Form.Get("details")
	todoToAdd := newTodo(name, details)

	err := todoHTML.ExecuteTemplate(w, "todo", todoToAdd)
	if err != nil {
		http.Error(w, "Failed to parse template", http.StatusBadRequest)
		return
	}
}

func reloadTodos(w http.ResponseWriter) {
	err := todoHTML.ExecuteTemplate(w, "todosList", BigBoy{
		Todos:   Todos,
		Todones: Todones,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
