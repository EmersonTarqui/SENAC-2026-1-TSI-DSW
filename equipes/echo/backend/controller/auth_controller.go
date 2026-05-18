package controller

import (
	"backend/model"
	"backend/repository"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var dadosDigitados model.User

	json.NewDecoder(r.Body).Decode(&dadosDigitados)

	usuarioDoBanco, err := repository.GetUserByEmail(dadosDigitados.Email)
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	// Compara a senha digitada com o hash que veio do banco
	err = bcrypt.CompareHashAndPassword([]byte(usuarioDoBanco.Password), []byte(dadosDigitados.Password))
	if err != nil {
		http.Error(w, "Credenciais inválidas", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Limpa a senha para não devolver o hash para o frontend
	usuarioDoBanco.Password = ""
	json.NewEncoder(w).Encode(usuarioDoBanco)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var novoUsuario model.User

	json.NewDecoder(r.Body).Decode(&novoUsuario)

	usuarioExistente, _ := repository.GetUserByEmail(novoUsuario.Email)

	if usuarioExistente.ID > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "E-mail já cadastrado no sistema"})
		return
	}

	// Gera o hash antes de salvar
	hash, err := bcrypt.GenerateFromPassword([]byte(novoUsuario.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Erro ao criptografar senha", http.StatusInternalServerError)
		return
	}

	// Substitui a senha em texto limpo pelo hash gerado
	novoUsuario.Password = string(hash)

	repository.CreateUser(novoUsuario)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuário registrado com sucesso"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// apenas avisamos o front para limpar tudo
	json.NewEncoder(w).Encode(model.Response{
		Status:  "success",
		Message: "Logout realizado com sucesso.",
	})
}
