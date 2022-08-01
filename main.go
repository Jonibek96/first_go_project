package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/gorilla/mux"
	"html/template"
	"net/http"
)

type User struct {
	Name                  string
	Age                   uint16
	Money                 int16
	Avq_grades, Happiness float64
	Hobbies               []string
}
type Articles struct {
	Id       uint16
	Title    string
	Anons    string
	Full_str string
}

var posts = []Articles{}
var showPost = Articles{}

func home_page(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golag")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	res, err := db.Query("SELECT * FROM `articles`")
	if err != nil {
		panic(err)
	}
	posts = []Articles{}
	for res.Next() {
		var post Articles
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_str)
		if err != nil {
			panic(err)
		}
		posts = append(posts, post)
	}

	tmpl.ExecuteTemplate(writer, "index", posts)
}

func create(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("templates/create.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}
	tmpl.ExecuteTemplate(writer, "create", nil)
}

func save_article(writer http.ResponseWriter, request *http.Request) {
	title := request.FormValue("title")
	anons := request.FormValue("anons")
	full_str := request.FormValue("full_str")
	if title == "" || anons == "" || full_str == "" {
		fmt.Fprintf(writer, "Error data")
	} else {
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golag")
		if err != nil {
			panic(err)
		}

		insert, err := db.Query(fmt.Sprintf("INSERT INTO  `articles` (`title`, `anons`, `full_str`) VALUES ( '%s', '%s', '%s')",
			title, anons, full_str))
		if err != nil {
			panic(err)
		}
		defer insert.Close()
		http.Redirect(writer, request, "/create/", http.StatusSeeOther)
	}
}

func handleRequest() {
	router := mux.NewRouter()
	router.HandleFunc("/", home_page).Methods("GET")
	router.HandleFunc("/contacts/", contact_page)
	router.HandleFunc("/create/", create).Methods("GET")
	router.HandleFunc("/save_article", save_article).Methods("POST")
	router.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", router)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":9080", nil)
}

func show_post(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golag")
	fmt.Println(err)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tmpl, err := template.ParseFiles("templates/show.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}
	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE  `id` = '%s'", vars["id"]))
	if err != nil {
		panic(err)
	}

	showPost = Articles{}
	for res.Next() {
		var post Articles
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.Full_str)
		if err != nil {
			panic(err)
		}
		showPost = post
	}

	tmpl.ExecuteTemplate(writer, "show", showPost)
}

func contact_page(writer http.ResponseWriter, request *http.Request) {
	tmpl, err := template.ParseFiles("templates/contacts.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(writer, err.Error())
	}
	tmpl.ExecuteTemplate(writer, "contacts", nil)
}

func myDataBase() {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golag")
	fmt.Println(err)
	if err != nil {
		panic(err)
	}

	insert, err := db.Query("INSERT INTO  `users` (`id`, `name`, `age`) VALUES (1, 'Alex', 25)")
	if err != nil {
		panic(err)
	}
	defer insert.Close()
	fmt.Println("ok mysql")

}

func main() {
	handleRequest()
}
