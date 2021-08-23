package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/tula3and/me-sign/email"
)

const (
	templateDir string = "templates/"
	port        string = ":4000"
)

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(rw, "home", nil)
}

func sign(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "sign", nil)
	case "POST":
		r.ParseForm()
		address := r.Form.Get("address")
		http.Redirect(rw, r, "/sent?email="+address, http.StatusPermanentRedirect)
	}
}

func sent(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		r.ParseForm()
		address := r.URL.Query().Get("email")
		verify := email.Verify(address)
		var data string
		if verify {
			data = "Success: Sent to " + address
		} else {
			data = "Failed: check your input again"
		}
		templates.ExecuteTemplate(rw, "sent", data)
	}
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/sign", sign)
	http.HandleFunc("/sent", sent)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
