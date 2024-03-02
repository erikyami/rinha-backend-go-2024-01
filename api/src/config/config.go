package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DB_USER  = ""
	DB_SENHA = ""
	DB_NOME  = ""
	DB_HOST  = ""
	DB_PORT  = 0
	API_PORT = 0
)

func Carregar() {
	var erro error
	if erro = godotenv.Load(); erro != nil {
		log.Fatal(erro)
	}

	API_PORT, erro = strconv.Atoi(os.Getenv("API_PORT"))
	if erro != nil {
		API_PORT = 5000
	}

	DB_PORT, erro = strconv.Atoi(os.Getenv("DB_PORT"))
	if erro != nil {
		DB_PORT = 5432
	}

	DB_USER = os.Getenv("DB_USUARIO")
	DB_SENHA = os.Getenv("DB_SENHA")
	DB_NOME = os.Getenv("DB_NOME")
	DB_HOST = os.Getenv("DB_HOST")
}
