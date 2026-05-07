package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()

	r := mux.NewRouter()

	r.HandleFunc("/tasks", getTasksHandler).Methods("GET")
	r.HandleFunc("/tasks", createTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{id}", updateTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTaskHandler).Methods("DELETE")

	http.ListenAndServe(":8080", enableCORS(r))
}