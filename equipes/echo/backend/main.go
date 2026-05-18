package main

import (
	"backend/controller"
	"backend/database"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Configurar corretamente os headers de segurança da API
		w.Header().Set("X-Content-Type-Options", "nosniff")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	database.Connect()

	defer database.DB.Close()

	r := mux.NewRouter()

	// tasks (v1)
	r.HandleFunc("/api/v1/tasks", controller.GetTasks).Methods("GET")
	r.HandleFunc("/api/v1/tasks", controller.CreateTask).Methods("POST")
	r.HandleFunc("/api/v1/tasks/{id}", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/v1/tasks/{id}", controller.DeleteTask).Methods("DELETE")

	// rota auth(users) (v1)
	r.HandleFunc("/api/v1/register", controller.Register).Methods("POST")
	r.HandleFunc("/api/v1/login", controller.Login).Methods("POST")
	r.HandleFunc("/api/v1/logout", controller.Logout).Methods("POST")

	// Melhor maneira: define o arquivo principal na raiz para nao listar diretorio
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/login.html")
	})

	// Esta rota deve vir por ultimo. Ela serve os arquivos da pasta frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("frontend/")))

	// Configura o servidor para permitir o graceful shutdown
	server := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: enableCORS(r),
	}

	// Canal para escutar sinais de parada do sistema
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Inicia o servidor em segundo plano para nao travar o canal
	go func() {
		log.Println("Servidor rodando na porta 8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro no servidor: %v", err)
		}
	}()

	// Aguarda o sinal do Podman/sistema para desligar
	<-stop
	log.Println("Desligando o servidor de forma segura (graceful shutdown)")

	// Tempo limite de 5 segundos para encerrar as conexoes ativas com seguranca
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
