package main

import (
	"encoding/json"
	"net/http"
	// "github.com/gorilla/mux"
	// "strconv"
	// "fmt"
	// "strings"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ================= VALIDAÇÃO =================

func validateTask(task Task, requireTitle bool) error {
	if requireTitle {
		if strings.TrimSpace(task.Title) == "" {
			return fmt.Errorf("title é obrigatório")
		}
		if len(task.Title) > 100 {
			return fmt.Errorf("title muito longo (máx 100 caracteres)")
		}
	}
	return nil
}

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