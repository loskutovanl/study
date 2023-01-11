package handlers

import (
	"30/internal/databaseRequests"
	"30/internal/entity"
	"30/internal/usecase"
	"30/internal/usecase/repo"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

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
		var u *entity.User
		if err = json.Unmarshal(content, &u); err != nil {
			log.Warn("Inside CreateHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// валидация данных пользователя, добавление пользователя в таблицу "users",
		// TODO to delete
		//userId, err := databaseRequests.ValidateUserAndCreateUser(db, u)

		// Use case
		userUseCase := usecase.New(
			repo.NewPostgreSQLClassicRepository(db),
		)

		userId, err := userUseCase.NewUser(u)

		//userId, err := userUsecase.NewUser(u)
		if err != nil {
			log.Errorf("Inside CreateHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully created user (user_id %d)", userId)

		// добавление связи друзей в таблицу "friends
		for _, friend := range u.Friends {
			err = databaseRequests.MakeFriendsForCreatedUser(db, friend, userId)
			if err != nil {
				log.Errorf("Inside CreateHandler: %s", err)
			} else {
				log.Infof("Successfully added friends relation (user1_id %d, user2_id %s) to database table friends", userId, friend)
			}
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(strconv.Itoa(userId)))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside CreateHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

// MakeFriendsHandler обрабатывает POST-запрос на дружбу двух пользователей.
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
		var f *entity.Friends
		if err = json.Unmarshal(content, &f); err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обработка и проверка значений пользовательских id
		err = databaseRequests.ValidateUsersIdAndMakeFriends(db, f.SourceId, f.TargetId)
		if err != nil {
			log.Errorf("Inside MakeFriendsHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		log.Infof("Successfully added friends relation (user1_id %s, user2_id %s) to database table friends", f.SourceId, f.TargetId)
		successMsg := fmt.Sprintf("%s и %s теперь друзья", f.SourceId, f.TargetId)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside MakeFriendsHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
