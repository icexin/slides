package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	db *sqlx.DB
)

type User struct {
	Id       int
	Name     string
	Password string
	Note     string
	Isadmin  bool
}

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

func Login(w http.ResponseWriter, r *http.Request) {
	render(w, "login", nil)
}

func CheckLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	user := r.FormValue("user")
	password := r.FormValue("password")

	var s string
	err := db.Get(&s, "SELECT password FROM user where name = ?", user)
	if err != nil {
		log.Print(err)
		http.Error(w, "bad user/password", 400)
		return
	}

	md5sum := md5.Sum([]byte(password))
	if s != fmt.Sprintf("%x", md5sum) {
		http.Error(w, "bad user/password", 400)
		return
	}

	cookie := http.Cookie{
		Name:   "user",
		Value:  user,
		Path:   "/",
		MaxAge: 600,
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/list", 302)
}

func List(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := db.Select(&users, "SELECT * FROM user")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	render(w, "list", users)
}

func Add(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.FormValue("name")
	passwd := fmt.Sprintf("%x", md5.Sum([]byte(r.FormValue("password"))))
	note := r.FormValue("note")
	_, err := db.Exec("INSERT INTO user VALUES(NULL,?,?,?,?)", name, passwd, note, 0)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.FormValue("id")
	passwd := fmt.Sprintf("%x", md5.Sum([]byte(r.FormValue("password"))))
	note := r.FormValue("note")
	_, err := db.Exec("UPDATE user SET password = ?, note=? WHERE id=?", passwd, note, id)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	idstr := r.FormValue("id")
	if idstr == "" {
		http.Error(w, "empty id", 400)
		return
	}
	id, err := strconv.Atoi(idstr)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	log.Printf("id:%v", id)
	_, err = db.Exec("DELETE FROM user WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func UploadPage(w http.ResponseWriter, r *http.Request) {
	render(w, "upload", nil)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		log.Print(err)
		return
	}
	f, _ := os.Create(filepath.Base(handler.Filename))
	defer f.Close()
	io.Copy(f, file)
}

func NeedLogin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("user")
		if err == nil {
			h(w, r)
			return
		}

		log.Print(err)
		http.Redirect(w, r, "/login", 302)
	}
}

func main() {
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", NeedLogin(List))
	http.HandleFunc("/login", Login)
	http.HandleFunc("/checkLogin", CheckLogin)
	http.HandleFunc("/list", NeedLogin(List))
	http.HandleFunc("/add", NeedLogin(Add))
	http.HandleFunc("/update", NeedLogin(Update))
	http.HandleFunc("/delete", NeedLogin(Delete))
	http.HandleFunc("/uploadPage", UploadPage)
	http.HandleFunc("/upload", Upload)

	var err error
	db, err = sqlx.Open("mysql", "golang:golang@tcp(59.110.12.72:3306)/go")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", nil))
}
