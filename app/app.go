package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/service"
	"github.com/rokafela/udemy-banking-lib/logger"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Panic("Error loading .env file")
	}

	mandatory_env := []string{
		"APP_ADDRESS",
		"APP_PORT",
		"DB_ADDRESS",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
	}

	for _, v := range mandatory_env {
		_, db_address_bool := os.LookupEnv(v)
		if !db_address_bool {
			logger.Panic("Environment variable not defined|" + v)
		}
	}
}

func createDbPool() *sqlx.DB {
	db_user := os.Getenv("DB_USER")
	db_password := os.Getenv("DB_PASSWORD")
	db_address := os.Getenv("DB_ADDRESS")
	db_port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")

	dbProperties := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_password, db_address, db_port, db_name)

	client, err := sqlx.Connect("mysql", dbProperties)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// See "Important settings" section.
	client.SetConnMaxLifetime(time.Minute * 3)
	client.SetMaxOpenConns(10)
	client.SetMaxIdleConns(10)

	return client
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Panic(err.Error())
	}
}

// var TestAuthHandler = AuthHandler{}

func Start() {
	loadEnv()

	// router
	router := mux.NewRouter()

	// repository initialization
	client := createDbPool()
	user_repository_db := domain.NewAuthRepositoryDb(client)

	// handler initialization
	auth_service := service.NewAuthService(user_repository_db)
	auth_handler := AuthHandler{auth_service}

	// routes
	router.HandleFunc("/login", auth_handler.HandleLogin).Methods(http.MethodPost)
	router.HandleFunc("/verify", auth_handler.HandleVerify).Methods(http.MethodPost)

	// server
	logger.Info(fmt.Sprintf("application listening in %s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")), router)
	if err != nil {
		logger.Panic(err.Error())
	}
}
