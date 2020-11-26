package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	tracer "github.com/neel229/tracing"
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
	r.tracer = tracer.New(os.Stdout)

	// Chat Path
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))

	// Login Path
	http.Handle("/login", &templateHandler{filename: "login.html"})

	// Start the room
	http.Handle("/room", r)
	go r.run()

	//Start the server and listen at port 8080
	fmt.Println("Starting server at port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
