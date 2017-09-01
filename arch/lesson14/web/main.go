package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	_ "github.com/go-sql-driver/mysql"
)

func render(w http.ResponseWriter, name string, data interface{}) {
	tplfile := filepath.Join("template", name+".tpl")
	tpl, err := template.ParseFiles(tplfile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "login", nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// fmt.Fprintf(w, "%s:%s", r.FormValue("user"), r.FormValue("password"))
	http.Redirect(w, r, "/list", 302)
}

func List(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "list")
}

func Add(w http.ResponseWriter, r *http.Request) {
}

func Update(w http.ResponseWriter, r *http.Request) {
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/list", List)
	db, err := sql.Open("mysql", "golang:golang@tcp(reboot:3306)/go")
	if err != nil {
		log.Fatal(err)
	}

	var user string
	row := db.QueryRow("SELECT CURRENT_USER()")
	err = row.Scan(&user)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("user:", user)
	//log.Fatal(http.ListenAndServe(":8080", nil))
}
