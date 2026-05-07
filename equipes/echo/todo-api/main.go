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

	http.ListenAndServe(":8080", enableCORS(r))
}