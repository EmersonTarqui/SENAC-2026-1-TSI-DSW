package controller

import (
	"backend/model"
	"backend/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func validateTask(task model.Task) error {
	if strings.TrimSpace(task.Title) == "" {
		return fmt.Errorf("title é obrigatório")
	}

	if len(task.Title) > 100 {
		return fmt.Errorf("title muito longo")
	}

	return nil
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Pegamos o ID do usuário que vem na URL (Ex: /tasks?user_id=1)
	// esse ID virá do seu sessionStorage no frontend.
	idUsuario := r.URL.Query().Get("user_id")
	userID, _ := strconv.Atoi(idUsuario)

	// Chamamos o banco passando o ID para filtrar
	tasks, err := repository.GetTasks(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	json.NewEncoder(w).Encode(model.Response{
		Status: "success",
		Data:   tasks,
	})
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var task model.Task

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: "JSON inválido",
		})

		return
	}

	err = validateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	err = repository.CreateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(model.Response{
		Status:  "success",
		Message: "tarefa criada",
		Data:    task,
	})
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: "ID inválido",
		})

		return
	}

	var task model.Task

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err = decoder.Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: "JSON inválido",
		})

		return
	}

	err = validateTask(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	err = repository.UpdateTask(id, task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	json.NewEncoder(w).Encode(model.Response{
		Status:  "success",
		Message: "tarefa atualizada",
	})
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: "ID inválido",
		})

		return
	}

	var task model.Task

	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: "JSON inválido",
		})

		return
	}

	err = repository.DeleteTask(id, task.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		json.NewEncoder(w).Encode(model.Response{
			Status:  "error",
			Message: err.Error(),
		})

		return
	}

	json.NewEncoder(w).Encode(model.Response{
		Status:  "success",
		Message: "tarefa deletada",
	})
}