package main

import (
    "database/sql"
    "log"
    "net/http"
    "text/template"
    _ "github.com/go-sql-driver/mysql"
)

// Full CRUD application in GoLang

// Setup #1 install MySQL make a new database create a table with 

// Employee structure that demonstrates the employee in the database
type Employee struct {
    Id int
    Name string
    City string
}

// dbInit Initializes a SQL DB for the web page
func dbInit () {
    db, err:= sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/")
    if err != nil{
        panic(err.Error())
    } else {
        _, err = db.Exec("CREATE DATABASE goblog")
        if err != nil {
            log.Println(err.Error())
        } else {
            log.Println("Database Created")
        }
        db.Exec("USE goblog")
        stmt, err := db.Prepare("CREATE TABLE `employee` (`id` int(6) unsigned NOT NULL AUTO_INCREMENT,`name` varchar(30) NOT NULL,`city` varchar(30) NOT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;")
        if err != nil {
            log.Println(err.Error())
        } else {
            _, err = stmt.Exec()
            if err != nil {
                log.Println(err.Error())
            } else {
                log.Println("Table Created")
            }
        }
    }
}

// dbConn returns a reference to the SQL database
func dbConn() (db *sql.DB) {
    dbDriver := "mysql"
    dbUser := "root"
    dbPass := "root"
    dbName := "goblog"


    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil {
        panic(err.Error())
    }
    return db
}

var tmpl = template.Must(template.ParseGlob("form/*"))

// Index pulls all the entries from the database and renders them to the Index.tmpl
func Index(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    selDB, err := db.Query("SELECT * FROM Employee ORDER BY id DESC")
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    res := []Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
        res = append(res, emp)
    }
    tmpl.ExecuteTemplate(w, "Index", res)
    defer db.Close()
}

// Show pulls one specific entry in the database and loads the Show.tmpl
func Show(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
    }
    tmpl.ExecuteTemplate(w, "Show", emp)
    defer db.Close()
}

// New loads the page for New.tmpl which is a form for inputting new entries into the DB
func New(w http.ResponseWriter, r *http.Request) {
    tmpl.ExecuteTemplate(w, "New", nil)
}

// Edit Renders the page for editing an existing entry and loads the text of the existing entry into the text firelds
func Edit(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    nId := r.URL.Query().Get("id")
    selDB, err := db.Query("SELECT * FROM Employee WHERE id=?", nId)
    if err != nil {
        panic(err.Error())
    }
    emp := Employee{}
    for selDB.Next() {
        var id int
        var name, city string
        err = selDB.Scan(&id, &name, &city)
        if err != nil {
            panic(err.Error())
        }
        emp.Id = id
        emp.Name = name
        emp.City = city
    }
    tmpl.ExecuteTemplate(w, "Edit", emp)
    defer db.Close()
}

// Insert Handles the database query where it inserts a new entry to the database
func Insert(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        city := r.FormValue("city")
        insForm, err := db.Prepare("INSERT INTO Employee(name, city) VALUES(?,?)")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, city)
        log.Println("INSERT: Name: " + name + " | City: " + city)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

// Update handles the database query where it updates an entry in the database
func Update(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    if r.Method == "POST" {
        name := r.FormValue("name")
        city := r.FormValue("city")
        id := r.FormValue("uid")
        insForm, err := db.Prepare("UPDATE Employee SET name=?, city=? WHERE id=?")
        if err != nil {
            panic(err.Error())
        }
        insForm.Exec(name, city, id)
        log.Println("UPDATE: Name: " + name + " | City: " + city)
    }
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

// Delete Handles the database query where it removes an entry fro mthe database 
func Delete(w http.ResponseWriter, r *http.Request) {
    db := dbConn()
    emp := r.URL.Query().Get("id")
    delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
    if err != nil {
        panic(err.Error())
    }
    delForm.Exec(emp)
    log.Println("DELETE")
    defer db.Close()
    http.Redirect(w, r, "/", 301)
}

func main() {
    // this is run once to Initialize the DB upon running 
    // it will log some errors stating the db already exists if youve ran the program before dont worry it keeps going if that happens
    dbInit()
    // URL routing
    log.Println("Server started on: http://localhost:8000")

    // Getters (GET)
    http.HandleFunc("/", Index) // READ or Get
    http.HandleFunc("/show", Show) // READ or Get
    http.HandleFunc("/new", New)
    http.HandleFunc("/edit", Edit) // READ or Get

    // setters (POST)
    http.HandleFunc("/insert", Insert) // CREATE or POST
    http.HandleFunc("/update", Update) // UPDATE or Put
    http.HandleFunc("/delete", Delete) // DELETE

    // setting up the server 
    http.ListenAndServe(":8000", nil)
}