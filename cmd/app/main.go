package main

import (
	"30/config"
	"30/internal/handlers"
	"30/internal/usecase"
	"30/internal/usecase/repo"
	"database/sql"
	"fmt"
	"net/http"
	"os"

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
	mux.Post("/create", func(w http.ResponseWriter, r *http.Request) { handlers.CreateHandler(w, r, db) })
	mux.Post("/make_friends", func(w http.ResponseWriter, r *http.Request) { handlers.MakeFriendsHandler(w, r, db) })
	mux.Delete("/user", func(w http.ResponseWriter, r *http.Request) { handlers.DeleteHandler(w, r, db) })
	mux.Get("/friends/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { handlers.GetAllFriendsHandler(w, r, db) })
	mux.Put("/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { handlers.PutUserAgeHandler(w, r, db) })

	err = http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Error("Unable to listen and serve:", err)
		return
	}
}
