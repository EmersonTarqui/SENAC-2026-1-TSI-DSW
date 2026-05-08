package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Connect() {
	user := "root"
	password := ""
	host := "localhost"
	port := "3306"
	dbname := "todo_db"

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		user,
		password,
		host,
		port,
		dbname,
	)

	var err error

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(
			"ERRO: Configuração do banco inválida! ",
			err,
		)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(
			"ERRO: Não consegui falar com o MySQL! ",
			err,
		)
	}

	fmt.Println("Conectado ao MySQL com sucesso!")
}