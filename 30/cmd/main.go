package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = "123456"
	dbname   = "server"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

type Storage struct {
	repository map[int]*User
}

func NewStorage() *Storage {
	return &Storage{
		make(map[int]*User),
	}
}

func (s *Storage) GetNextId() int {
	return len(s.repository) + 1
}

func (s *Storage) MakeFriends(sourceId, targetId int) (string, error) {
	var sourceUser *User
	var targetUser *User

	for id, user := range s.repository {
		switch id {
		case sourceId:
			sourceUser = user
		case targetId:
			targetUser = user
		}
	}

	if sourceUser == nil {
		return "", fmt.Errorf("no *User found with given sourceId %d", sourceId)
	}
	if targetUser == nil {
		return "", fmt.Errorf("no *User found with given targetId %d", targetId)
	}

	//s.repository[sourceId].Friends = append(s.repository[sourceId].Friends, targetUser)
	//s.repository[targetId].Friends = append(s.repository[targetId].Friends, sourceUser)
	return fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name), nil
}

func (s *Storage) DeleteUser(targetId int) (string, error) {
	var userToDelete *User

	for id, user := range s.repository {
		if id == targetId {
			userToDelete = user
		}
	}

	if userToDelete == nil {
		return "", fmt.Errorf("no *User found with given targetId %d", targetId)
	}

	//for _, friend := range userToDelete.Friends {
	//	for id, user := range s.repository {
	//		if user == friend {
	//			delete(s.repository, id)
	//			break
	//		}
	//	}
	//}
	delete(s.repository, targetId)
	return userToDelete.Name, nil
}

// User содержит информацию о пользователе: имя, возраст, список друзей
type User struct {
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

// Friends содержит информацию об id двух пользователей, отправивших запрос на дружбу,
type Friends struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

type DeleteUser struct {
	TargetId string `json:"target_id"`
}

func main() {
	// подключение, открытие и отложенное закрытие базы данных
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

	storage := NewStorage()

	// создание роутера и регистрация хендлеров
	r := chi.NewRouter()
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) { CreateHandler(w, r, db) })
	r.Post("/make_friends", func(w http.ResponseWriter, r *http.Request) { MakeFriendsHandler(w, r, db) })
	r.Delete("/user", storage.DeleteHandler)

	err = http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Error("Unable to listen and serve:", err)
		return
	}
}

// CreateHandler обрабатывает POST-запрос на добавление нового пользователя в базу данных.
// CreateHandler принимает три параметра: экземпляры райтера, респонса и базы данных.
// Обрабатывает ошибки подключения и преобразования данных, демаршаллизирует json-запрос. Добавляет в таблицу "users"
// имя и возраст нового пользователя, добавляет (при наличии) в таблицу "friends" id нового добавленного пользователя
// и всех его друзей. Логирует возможные ошибки. При успешном запросе возвращает ID пользователя и статус 201.
// При неуспешном запросе возвращает статус 400 (неверный метод) ии 500 (ошибка обработки данных).
func CreateHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Info("Inside CreateHandler")

	if r.Method == "POST" {
		// чтение запроса и обработка ошибок
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside CreateHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// демаршализация запроса и обработка ошибок
		var u *User
		if err = json.Unmarshal(content, &u); err != nil {
			log.Warn("Inside CreateHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// приведение типов возраста пользователя, запись имени и возраста пользователя в таблицу "users"
		age, err := strconv.Atoi(u.Age)
		if err != nil {
			log.Errorf("Unable to convert age %s from string to int type: %s", u.Age, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		lastInsertId := 0
		err = db.QueryRow("INSERT INTO users (name, age) VALUES($1, $2) RETURNING id", u.Name, age).Scan(&lastInsertId)
		if err != nil {
			log.Errorf("Unable to insert user (name %s, age %d) to database table users: %s", u.Name, age, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully added user (id %d) to database", lastInsertId)

		// приведение типов друзей пользователя, запись друзей в таблицу "friends"
		for _, friend := range u.Friends {
			friendId, err := strconv.Atoi(friend)
			if err != nil {
				log.Errorf("Unable to convert friendId %s to int: %s", friend, err)
			}
			_, err = db.Exec(`insert into "friends"("user1_id", "user2_id") values($1, $2)`, lastInsertId, friend)
			if err != nil {
				log.Errorf("Unable to insert friends (user1_id %d, user2_id %d) to database table friends: %s", lastInsertId, friendId, err)
			}
			log.Infof("Successfully added friends relation (user1_id %d, user2_id %d) to database table friends", lastInsertId, friendId)
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(strconv.Itoa(lastInsertId)))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside CreateHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

// MakeFriendsHandler обрабатывает POST-запрос на дружьу двух пользователей.
// MakeFriendsHandler принимает три параметра: экземпляры райтера, респонса и базы данных.
// Обрабатывает ошибки подключения и преобразования данных, демаршаллизирует json-запрос. Добавляет в таблицу
// "friends" id двух пользователей. Логирует возможные ошибки. При успешном запросе возвращает ID пользователя
// и статус 200. При неуспешном запросе возвращает статус 400 (неверный метод) ии 500 (ошибка обработки данных).
func MakeFriendsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Info("Inside MakeFriendsHandler")

	// чтение запроса и обработка ошибок
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// демаршализация запроса и обработка ошибок
		var f *Friends
		if err = json.Unmarshal(content, &f); err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// приведение типов друзей пользователя
		sourceId, err := strconv.Atoi(f.SourceId)
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to convert sourceId %d to int: %s", sourceId, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		targetId, err := strconv.Atoi(f.TargetId)
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to convert targetId %d to int: %s", targetId, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// проверка, что указанные id пользователей существуют в таблице "users"
		rows, err := db.Query(`select "id" from "users" where "id" = $1 or "id" = $2`, sourceId, targetId)
		defer rows.Close()

		if err != nil {
			log.Warnf("Inside MakeFrindsHandler, unable to perform select query on users table in database: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		var users []int
		var iuser int
		for rows.Next() {
			rows.Scan(&iuser)
			users = append(users, iuser)
		}
		if len(users) != 2 {
			log.Warnf("Inside MakeFriendsHandler, unable to find users with targetId %d and/or sourceI %d in users table in database: %s", targetId, sourceId, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// проверка, что указанные id пользователей еще не являются друзьями

		// запись друзей в таблицу "friends"
		_, err = db.Exec(`insert into "friends"("user1_id", "user2_id") values($1, $2)`, sourceId, targetId)
		if err != nil {
			log.Errorf("Unable to insert friends (user1_id %d, user2_id %d) to database table friends: %s", sourceId, targetId, err)
		}
		log.Infof("Successfully added friends relation (user1_id %d, user2_id %d) to database table friends", sourceId, targetId)

		successMsg := fmt.Sprintf("%d и %d теперь друзья", sourceId, targetId)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside MakeFriendsHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Storage) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside DeleteHandler")
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside DeleteHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				log.Warn("Inside DeleteHandler, unable to write response after reading http.Request.Body:", err2)
				return
			}
			return
		}

		var du *DeleteUser
		if err = json.Unmarshal(content, &du); err != nil {
			log.Warn("Inside DeleteHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		targetId, err := strconv.Atoi(du.TargetId)
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable to convert targetId %d to int: %s", targetId, err)
			return
		}

		requestStatus, err := s.DeleteUser(targetId)
		fmt.Println(s.repository)
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable delete user %d: %s", targetId, err)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(requestStatus))
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable to write response afterdeleting user %d: %s", targetId, err)
		}
		return
	}

	log.Infof("Inside DeleteHandler, inappropriate http.Request.Method: DELETE required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
