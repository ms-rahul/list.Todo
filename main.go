package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

var db *sql.DB

var rows []Row
var temp *template.Template

func init() {
	temp = template.Must(template.ParseGlob("templates/*.html"))
}

func startServer() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/add-new-todo", postToDo)
	http.HandleFunc("/current-todo-list", getToDo)
	http.HandleFunc("/get-single-todo", getById)
	http.HandleFunc("/update-existing-todo", putToDo)
	http.HandleFunc("/delete-from-todo", deleteToDo)

	port := importFromEnv("PORT")
	err := http.ListenAndServe(":"+port, nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err == nil {
		fmt.Println("server started at port:", port)
	} else {
		fmt.Println("error starting server: ", err)
		os.Exit(1)
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")

	temp.ExecuteTemplate(w, "root.html", nil)
}

func postToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /add-new-todo request")

	temp.ExecuteTemplate(w, "add-task.html", nil)
	desc := r.FormValue("todo_desc")
	status := r.FormValue("status")
	submitValue := r.FormValue("submit")
	if len(submitValue) != 0 {
		if len(desc) > 0 {
			addRow(db, desc, status)
			fmt.Fprintf(w, "Your To-Do entry has been added!")
		} else {
			fmt.Fprintf(w, "Task Description cannot be empty!")
		}
	}
}

func getToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /current-todo-list request")

	getAllRows(db)
	temp.ExecuteTemplate(w, "todo-list.html", rows)
	rows = []Row{}
}

func getById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /get-single-todo request")

	temp.ExecuteTemplate(w, "one-task.html", nil)
	id := r.FormValue("id")
	submitValue := r.FormValue("submit")
	if len(submitValue) != 0 {
		if len(id) > 0 {
			desc, s := getSingleToDo(db, id)
			if desc == "" {
				fmt.Fprintf(w, "Enter a valid To-Do ID")
			} else {
				var status string
				if s == true {
					status = "Completed"
				} else {
					status = "Not Completed"
				}
				fmt.Fprintf(w, "The requested To-Do: %s | Status: %s", desc, status)
			}
		} else {
			fmt.Fprintf(w, "Enter a valid To-Do ID")
		}
	}
}

func putToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /update-existing-todo request")

	temp.ExecuteTemplate(w, "update-task.html", nil)
	id := r.FormValue("id")
	desc := r.FormValue("todo_desc")
	status := r.FormValue("status")
	submitValue := r.FormValue("submit")
	if len(submitValue) != 0 {
		if len(id) > 0 && len(desc) > 0 {
			check, reason := updateRow(db, id, desc, status)
			if check == false {
				if reason == "ID" {
					fmt.Fprintf(w, "Please check if entered ID is valid: %s", id)
				} else if reason == "equal" {
					fmt.Fprintf(w, "Entered To-Do entry matches exisiting To-Do entry, Please change the entry and try again.")
				} else {
					fmt.Fprintf(w, "Please check entered To-Do details and try again.")
				}
			} else {
				fmt.Fprintf(w, "The requested To-Do entry has been updated!")
			}
		} else {
			fmt.Fprintf(w, "Please check entered To-Do details and try again.")
		}
	}
}

func deleteToDo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("got /delete-from-todo request")

	temp.ExecuteTemplate(w, "delete-task.html", nil)
	id := r.FormValue("id")
	submitValue := r.FormValue("submit")

	if len(submitValue) != 0 {
		if len(id) > 0 {
			check := deleteRow(db, id)

			if check == true {
				fmt.Fprintf(w, "The required To-Do entry has been DELETED successfully!")
			} else {
				fmt.Fprintf(w, "The entered ID does not exist.")
			}
		} else {
			fmt.Fprintf(w, "Enter a Valid ID")
		}
	}
}

func main() {
	//creating a database to organize to-do data
	db = createAndUseDB()

	//creates file server to serve misc files on templates
	fs := http.FileServer(http.Dir("misc"))
	http.Handle("/misc/", http.StripPrefix("/misc", fs))

	//creating a handler to access database
	// bh := newHandler(db)

	//starting server on local machine
	startServer()
}
