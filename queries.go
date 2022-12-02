package main

import (
	"database/sql"
	"log"
	"strconv"
)

type ToDo struct {
	ID   int    `json:"ID"`
	DESC string `json:"Description"`
	STAT bool   `json:"Status"`
}

type Row struct {
	ID   int
	DESC string
	STAT string
}

func addRow(db *sql.DB, desc string, s string) {
	var status bool
	if s == "on" {
		status = true
	} else {
		status = false
	}
	r := ToDo{DESC: desc, STAT: status}

	query, err := db.Prepare("INSERT INTO todo_list (todo_desc, status) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := query.Exec(r.DESC, r.STAT)
	if err != nil {
		log.Fatal(err)
	}

	rows_aff, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d todo created ", rows_aff)
}

func getAllRows(db *sql.DB) {
	todoList := []ToDo{}
	result, err := db.Query("SELECT * FROM todo_list")
	if err != nil {
		log.Fatal(err)
	}

	for result.Next() {
		var r ToDo

		err = result.Scan(&r.ID, &r.DESC, &r.STAT)
		if err != nil {
			log.Fatal(err)
		}
		todoList = append(todoList, r)
	}
	for i := 0; i < len(todoList); i++ {
		id := todoList[i].ID
		desc := todoList[i].DESC
		s := todoList[i].STAT
		var row Row
		if s == true {
			row = Row{ID: id, DESC: desc, STAT: "Completed"}
		} else {
			row = Row{ID: id, DESC: desc, STAT: "Not Completed"}
		}
		rows = append(rows, row)

	}
}

func getSingleToDo(db *sql.DB, i string) (string, bool) {
	id, _ := strconv.Atoi(i)
	result, err := db.Query("SELECT * FROM todo_list WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	if result.Next() {
		var r ToDo
		err := result.Scan(&r.ID, &r.DESC, &r.STAT)

		if err != nil {
			log.Fatal(err)
		}
		return r.DESC, r.STAT
	} else {
		return "", false
	}
}

func updateRow(db *sql.DB, i string, desc string, s string) (bool, string) {
	id, _ := strconv.Atoi(i)
	ret_desc, ret_stat := getSingleToDo(db, i)

	if ret_desc == "" {
		return false, "ID"
	}

	var status bool
	if s == "on" {
		status = true
	} else {
		status = false
	}

	r := ToDo{ID: id, DESC: desc, STAT: status}

	if ret_desc == r.DESC && ret_stat == r.STAT {
		return false, "equal"
	}

	query, err := db.Prepare("UPDATE todo_list SET todo_desc=?, status=? WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := query.Exec(r.DESC, r.STAT, r.ID)
	if err != nil {
		log.Fatal(err)
	}

	rows_aff, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%d todo updated ", rows_aff)
	if rows_aff == 0 {
		return false, ""
	}

	return true, "ok"
}

func deleteRow(db *sql.DB, i string) bool {
	id, _ := strconv.Atoi(i)
	r := ToDo{ID: id}
	query, err := db.Prepare("DELETE FROM todo_list WHERE id=?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := query.Exec(r.ID)
	if err != nil {
		log.Fatal(err)
	}

	rows_aff, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%d todo deleted ", rows_aff)

	if rows_aff == 0 {
		return false
	} else {
		return true
	}
}
