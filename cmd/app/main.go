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
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx"
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

	// migrate -database "postgres://postgres:123456@localhost:5432/server?sslmode=disable" -path migrations up
	//m, err := migrate.New("file:///migrations", "postgres://postgres:123456@localhost:5432/server?sslmode=disable")
	//fmt.Println(err)
	//err = m.Up()
	//defer m.Close()

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
