package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/rokafela/udemy-banking-auth/domain"
	"github.com/rokafela/udemy-banking-auth/logger"
	"github.com/rokafela/udemy-banking-auth/service"
)

var validate *validator.Validate

func init() {
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

func Start() {
	// router
	router := mux.NewRouter()

	// repository initialization
	client := createDbPool()
	userRepositoryDb := domain.NewAuthRepositoryDb(client)

	// validator initialization
	english := en.New()
	universal_translator := ut.New(english, english)
	translator, found := universal_translator.GetTranslator("en")
	if !found {
		logger.Fatal("translator not found")
	}
	validate = validator.New()
	register_err := en_translations.RegisterDefaultTranslations(validate, translator)
	if register_err != nil {
		logger.Fatal(register_err.Error())
	}
	_ = validate.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	// handler initialization
	ah := AuthHandler{service.NewAuthService(userRepositoryDb), validate}

	// routes
	router.HandleFunc("/login", ah.HandleLogin).Methods(http.MethodPost)
	router.HandleFunc("/verify", ah.HandleVerify).Methods(http.MethodPost)

	// server
	logger.Info(fmt.Sprintf("application listening in %s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")))
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", os.Getenv("APP_ADDRESS"), os.Getenv("APP_PORT")), router)
	if err != nil {
		logger.Panic(err.Error())
	}
}
