package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/lit"
	_ "github.com/go-sql-driver/mysql"
)

// Full CRUD application in GoLang

// Setup #1 install MySQL make a new database create a table with

// Employee structure that demonstrates the employee in the database
type Employee struct {
	Id   int
	Name string
	City string
}

// Credentials structure that represents the connection information to input into the terminal
type Credentials struct {
	UserName string
	Password string
	DBName   string
}

func credChecker() Credentials {
	var creds Credentials
	_, err := os.Stat("configs/creds.json")
	if err != nil {
		lit.Debug("Create file route of credChecker")
		creds = getCreds()
		os.MkdirAll("configs", 0777)
		file, _ := json.MarshalIndent(creds, "", "\t")
		_ = ioutil.WriteFile("configs/creds.json", file, 0644)
	} else {
		// pulls the creds folder if there was no error on checking the status.
		lit.Debug("Read file route of credChecker")
		file, _ := ioutil.ReadFile("configs/creds.json")
		_ = json.Unmarshal(file, &creds)
	}
	return creds
}

func getCreds() Credentials {
	creds := Credentials{}
	db, _ := sql.Open("mysql", creds.UserName+":"+creds.Password+"@tcp(127.0.0.1:3306)/")
	validCred := db.Ping()

	for validCred != nil {
		reader := bufio.NewReader(os.Stdin)
		lit.Debug("getCreds post reader")
		// prompts user for UserName
		fmt.Println("User Name:")
		creds.UserName, _ = reader.ReadString('\n')
		creds.UserName = strings.Replace(creds.UserName, "\r\n", "", -1)
	
		// prompts user for Password
		fmt.Println("Password:")
		creds.Password, _ = reader.ReadString('\n')
		creds.Password = strings.Replace(creds.Password, "\r\n", "", -1)
	
		// prompts user for db name
		fmt.Println("Database Name:")
		creds.DBName, _ = reader.ReadString('\n')
		creds.DBName = strings.Replace(creds.DBName, "\r\n", "", -1)
		db, _ = sql.Open("mysql", creds.UserName+":"+creds.Password+"@tcp(127.0.0.1:3306)/")
		validCred = db.Ping()
		log.Println(validCred)
	}

	return creds
}

var creds Credentials = credChecker()

// createFile makes a file it takes 3 parameters
// path (string) the child directory from the root that the file goes into DO NOT MAKE THIS A BLANK STRING
// fileName (string) the name of the file and file type your writing
// panic (bool) triggers panic if there is an error or continues to run and just logs the error
func createFile(path string, fileName string, panicTrigger bool, content string) {
	_, err := os.Stat(path + "/" + fileName)
	if os.IsNotExist(err) {
		file, err := os.Create(path + "/" + fileName)
		defer file.Close()
		if err != nil {
			if panicTrigger {
				panic(err.Error())
			} else {
				log.Println(err.Error())
			}
		} else {
			file, err := os.OpenFile(path+"/"+fileName, os.O_RDWR, 0644)
			defer file.Close()
			if err != nil {
				if panicTrigger {
					panic(err.Error())
				} else {
					log.Println(err.Error())
				}
			}
			_, err = file.WriteString(content)
			if err != nil {
				if panicTrigger {
					panic(err.Error())
				} else {
					log.Panicln(err.Error())
				}
			}
			err = file.Sync()
			if err != nil {
				if panicTrigger {
					panic(err.Error())
				} else {
					log.Panicln(err.Error())
				}
			}
		}
	}
}

// creating templates from thin air Note: programmed in a sleep deprived detrmination
func createFileStructure() {
	lit.Debug("Enter Create File Structure")
	path := "form"
	os.MkdirAll(path, 0777)
	createFile(path, "Header.tmpl", false, "{{ define \"Header\" }} \n<!DOCTYPE html> \n<html lang=\"en-US\"> \n\t<head> \n\t\t<title>Golang MySQL CRUD Spawner</title> \n\t\t<meta charset=\"UTF-8\" /> \n\t</head> \n\t<body> \n\t\t<h1>Golang MySQL CRUD Spawner</h1> \n{{ end }}")
	createFile(path, "Edit.tmpl", false, "{{ define \"Edit\" }}\n\t{{ template \"Header\" }}\n\t\t{{ template \"Menu\" }}\n\t\t<h2>Edit Name and City</h2>\n\t\t<form method=\"POST\" action=\"update\">\n\t\t\t<input type=\"hidden\" name=\"uid\" value=\"{{ .Id }}\" />\n\t\t\t<label> Name </label><input type=\"text\" name=\"name\" value=\"{{ .Name }}\" /><br />\n\t\t\t<label> City </label><input type=\"text\" name=\"city\" value=\"{{ .City }}\" /><br />\n\t\t\t<input type=\"submit\" value=\"Save user\" />\n\t\t</form><br />\n\t{{ template \"Footer\" }}\n{{ end }}")
	createFile(path, "Footer.tmpl", false, "{{ define \"Footer\" }}\n</body>\n\n</html>\n{{ end }}")
	createFile(path, "Index.tmpl", false, "{{ define \"Index\" }}\n\t{{ template \"Header\" }}\n\t\t{{ template \"Menu\"  }}\n\t\t<h2> Registered </h2>\n\t\t<table border=\"1\">\n\t\t\t<thead>\n\t\t\t\t<tr>\n\t\t\t\t\t<td>ID</td>\n\t\t\t\t\t<td>Name</td>\n\t\t\t\t\t<td>City</td>\n\t\t\t\t\t<td>View</td>\n\t\t\t\t\t<td>Edit</td>\n\t\t\t\t\t<td>Delete</td>\n\t\t\t\t</tr>\n\t\t\t</thead>\n\t\t\t<tbody>\n\t\t\t\t{{ range . }}\n\t\t\t\t<tr>\n\t\t\t\t\t<td>{{ .Id }}</td>\n\t\t\t\t\t<td> {{ .Name }} </td>\n\t\t\t\t\t<td>{{ .City }} </td>\n\t\t\t\t\t<td><a href=\"/show?id={{ .Id }}\">View</a></td>\n\t\t\t\t\t<td><a href=\"/edit?id={{ .Id }}\">Edit</a></td>\n\t\t\t\t\t<td><a href=\"/delete?id={{ .Id }}\">Delete</a></td>\n\t\t\t\t</tr>\n\t\t\t\t{{ end }}\n\t\t\t</tbody>\n\t\t</table>\n\t{{ template \"Footer\" }}\n{{ end }}")
	createFile(path, "Menu.tmpl", false, "{{ define \"Menu\" }}\n<a href=\"/\">HOME</a> | \n<a href=\"/new\">NEW</a>\n{{ end }}")
	createFile(path, "New.tmpl", false, "call New.tmpl file inside form.\n\n{{ define \"New\" }}\n\t{{ template \"Header\" }}\n\t\t{{ template \"Menu\" }}\n\t\t<h2>New Name and City</h2>\n\t\t<form method=\"POST\" action=\"insert\">\n\t\t\t<label> Name </label><input type=\"text\" name=\"name\" /><br />\n\t\t\t<label> City </label><input type=\"text\" name=\"city\" /><br />\n\t\t\t<input type=\"submit\" value=\"Save user\" />\n\t\t</form>\n\t{{ template \"Footer\" }}\n{{ end }}")
	createFile(path, "Show.tmpl", false, "{{ define \"Show\" }}\n\t{{ template \"Header\" }}\n\t\t{{ template \"Menu\"  }}\n\t\t<h2> Register {{ .Id }} </h2>\n\t\t\t<p>Name: {{ .Name }}</p>\n\t\t\t<p>City:  {{ .City }}</p><br /> <a href=\"/edit?id={{ .Id }}\">Edit</a></p>\n\t{{ template \"Footer\" }}\n{{ end }}")

}

// dbInit Initializes a SQL DB for the web page
func dbInit() {
	//User Input
	usernm := creds.UserName
	pass := creds.Password
	DBName := creds.DBName
	lit.Debug("Hit dbInit " + DBName)
	log.Println(usernm + ":" + pass + "@tcp(127.0.0.1:3306)/")

	db, err := sql.Open("mysql", usernm+":"+pass+"@tcp(127.0.0.1:3306)/")
	err = db.Ping() //Need to ping to generate connection and trigger err
	if err != nil {
		lit.Error("Error in Init Log-in")
		creds = getCreds()
		file, _ := json.MarshalIndent(creds, "", "\t")
		_ = ioutil.WriteFile("configs/creds.json", file, 0644)
	} else {
		lit.Debug("Attempt DB Creation")
		_, err = db.Exec("CREATE DATABASE " + DBName)
		if err != nil {
			log.Println(err.Error())
		} else {
			log.Println("Database Created:", "\""+DBName+"\"")
		}
		db.Exec("USE " + DBName)
		stmt, err := db.Prepare("CREATE TABLE `employee` (`id` int(6) unsigned NOT NULL AUTO_INCREMENT,`name` varchar(30) NOT NULL,`city` varchar(30) NOT NULL,PRIMARY KEY (`id`));")
		if err != nil {
			log.Println(err.Error())
		} else {
			_, err = stmt.Exec()
			if err != nil {
				log.Println(err.Error())
			} else {
				log.Println("Table Created", "\""+"employees"+"\"")
			}
		}
	}
}

// dbConn returns a reference to the SQL database
func dbConn() (db *sql.DB) {
	// declaring some variables to use in the in the when opening the db not super useful now but later version will be helpful
	dbDriver := "mysql"
	dbUser := creds.UserName
	dbPass := creds.Password
	DBName := creds.DBName
	lit.Debug("Hit dbConn Entered")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+DBName)
	if err != nil {
		lit.Error("Something terrible has happened")
	}
	return db
}

// getTemplates this imports the templates from the form folder
// this is done for flow reasons upon startup as wel as allowing for the ability
// to make changes to the html formatting without restarting the server
func getTemplates() *template.Template {
	return template.Must(template.ParseGlob("form/*"))
}

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
	getTemplates().ExecuteTemplate(w, "Index", res)
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
	getTemplates().ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

// New loads the page for New.tmpl which is a form for inputting new entries into the DB
func New(w http.ResponseWriter, r *http.Request) {
	getTemplates().ExecuteTemplate(w, "New", nil)
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
	getTemplates().ExecuteTemplate(w, "Edit", emp)
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
	log.Println(emp)
	delForm, err := db.Prepare("DELETE FROM Employee WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delForm.Exec(emp)
	log.Println("DELETE")
	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func tester() {

}

func init() {
	lit.LogLevel = 3
	lit.Debug("Logging Level is Debug")

	credChecker()
	// this is run once to Initialize the DB upon running
	// it will log some errors stating the db already exists if youve ran the program before dont worry it keeps going if that happens
	dbInit()

	createFileStructure()

}

func main() {

	// URL routing
	log.Println("Server started on: http://localhost:8000")

	// Getters (GET)
	http.HandleFunc("/", Index)    // READ or Get
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
