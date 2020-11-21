package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once
	filename string
	temp     *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.temp = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	if err := t.temp.Execute(w, nil); err != nil {
		log.Fatal("Error executing the Template provided")
	}
}

func main() {
	// Create a new room
	r := newRoom()

	// Root Path
	http.Handle("/", &templateHandler{filename: "chat.html"})

	// Start the room
	http.Handle("/room", r)
	go r.run()

	//Start the server and listen at port 8080
	fmt.Println("Starting server at port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
