package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"study/config"
	"study/internal/controller/http/v1"
	"study/internal/usecase"
	"study/internal/usecase/repo"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// настройка логов и загрузка переменных окружения
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found: %s", err)
		os.Exit(1)
	}
}

func main() {
	// загрузка переменных, подключение и отложенное закрытие базы данных
	conf := config.New()
	var (
		host     = conf.Host
		port     = conf.Port
		password = conf.Password
		dbname   = conf.Dbname
		user     = conf.User
	)

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Error("Unable to open database:", err)
	}
	//driver, err := postgres.WithInstance(db, &postgres.Config{})
	//m, err := migrate.NewWithDatabaseInstance(
	//	"file:///migrations",
	//	"postgres", driver)
	//m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run

	defer func() {
		err = db.Close()
		if err != nil {
			log.Error("Unable to close database:", err)
		}
	}()

	// Use case
	userUseCase := usecase.New(
		repo.NewPostgreSQLClassicRepository(db),
	)

	// создание роутера и регистрация хендлеров
	mux := chi.NewRouter()
	v1.NewUserRoutes(mux, userUseCase)

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Error("Unable to listen and serve:", err)
		return
	}
}
