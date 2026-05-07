package main

import (
	"encoding/json"
	"net/http"
	// "github.com/gorilla/mux"
	// "strconv"
	// "fmt"
	// "strings"
)

// ================= ROTAS =================

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" {
		http.Error(w, "Método inválido", 405)
		return
	}

	getTasks(w)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "POST" {
		http.Error(w, "Método inválido", 405)
		return
	}

	createTask(w, r)
}

// ================= FUNCOES PRINCIPAIS =================

// ================= READ =================

func getTasks(w http.ResponseWriter) {
	rows, err := db.Query(`
		SELECT id, user_id, title, done, created_at
		FROM tasks
	`)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var tasks []Task

	for rows.Next() {
		var t Task

		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Title,
			&t.Done,
			&t.CreatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		tasks = append(tasks, t)
	}

	json.NewEncoder(w).Encode(tasks)
}

// ================= CREATE =================

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&task)
	if err != nil {
		http.Error(w, "JSON inválido", 400)
		return
	}

	if task.UserID == 0 {
		http.Error(w, "user_id é obrigatório", 400)
		return
	}

	err = validateTask(task, true)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	result, err := db.Exec(
		`INSERT INTO tasks(user_id, title, done)
		 VALUES(?, ?, ?)`,
		task.UserID,
		task.Title,
		false,
	)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	id, _ := result.LastInsertId()
	task.ID = int(id)
	task.Done = false

	json.NewEncoder(w).Encode(task)
}