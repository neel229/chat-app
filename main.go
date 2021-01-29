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
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

const (
	key    string = "725697710987-pj9hmpg8q7cls962skjbjs0v6v12diuj.apps.googleusercontent.com"
	secret string = "qrtOYbS--627XhkM0CSEZ5Me"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	if err := t.temp.Execute(w, data); err != nil {
		log.Fatalf("There was an error executing the template: %v", err)
	}
}

func main() {
	// OAuth setup
	gomniauth.SetSecurityKey("GcRcm0HhPeCjUK9Kdahy9rnYNxJ3olhDOPZXnFnfb2Y3NmFpE1NAoQvl6sZ9GGzf")
	gomniauth.WithProviders(
		google.New(key, secret, "http://localhost:8080/auth/callback/google"),
	)

	// Create a new room
	r := newRoom()
	r.tracer = tracer.New(os.Stdout)

	// Chat Path
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))

	// Login Path
	http.Handle("/login", &templateHandler{filename: "login.html"})

	// Auth
	http.HandleFunc("/auth/", loginHandler)
	// Start the room
	http.Handle("/room", r)
	go r.run()

	//Start the server and listen at port 8080
	fmt.Println("Starting server at port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
