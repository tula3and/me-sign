package main

import (
	"fmt"
	"html/template"
	"image/png"
	"log"
	"net/http"

	"github.com/tula3and/me-sign/email"
	"github.com/tula3and/me-sign/sign"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	templateDir string = "templates/"
	port        string = ":4000"
)

var templates *template.Template

func home(rw http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(rw, "home", nil)
}

func make(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "make", nil)
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
			data = "Success: sent to " + address
		} else {
			data = "Failed: check your input again"
		}
		templates.ExecuteTemplate(rw, "sent", data)
	}
}

func key(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		r.ParseForm()
		email := r.URL.Query().Get("email")
		signed := r.URL.Query().Get("signed")
		verify := sign.Verify(signed, email, sign.RestorePublicKey(sign.Key()))
		if verify {
			templates.ExecuteTemplate(rw, "realSign", nil)
		} else {
			http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
		}
	}
}

//create qrcode
func qrcodeviewCodeHandler(w http.ResponseWriter, r *http.Request) {
	dataString := r.FormValue("fileName")

	qrCode, _ := qr.Encode(dataString, qr.L, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 512, 512)

	png.Encode(w, qrCode)
}

func main() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))
	http.HandleFunc("/", home)
	http.HandleFunc("/make", make)
	http.HandleFunc("/sent", sent)
	http.HandleFunc("/key", key)
	fmt.Printf("Listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
