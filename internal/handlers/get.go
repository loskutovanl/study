package handlers

import (
	"30/internal/databaseRequests"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

// GetAllFriendsHandler обрабатывает GET-запрос на получение всех друзей пользователя. Логирует возможные ошибки.
// При успешном запросе возвращает список друзей пользователя (или сообщение об их отсутствии) и статус 200.
// При неуспешном запросе возвращает статус 400 (неверный метод) ии 500 (ошибка обработки данных).
func GetAllFriendsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Info("Inside GetAllFriendsHandler")

	// чтение запроса, приведение userId к числовому типу и обработка ошибок
	if r.Method == "GET" {
		userIdString := chi.URLParam(r, "id")
		userIdInt, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Warnf("Inside GetAllFriendsHandler, unable to convert user_id %s from string to int: %s", userIdString, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// проверка, что пользователь существует в таблице "users"
		err = databaseRequests.CheckUserExistsInUsersTable(db, userIdInt)
		if err != nil {
			log.Warnf("Inside GetAllFriendsHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// извлечение друзей пользователя из таблиц "users" и "friends"
		friends, err := databaseRequests.GetFriendsForUser(db, userIdInt)
		if err != nil {
			log.Errorf("Inside GetAllFriendsHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// запись сообщения
		var successMsg string
		if len(friends) == 0 {
			successMsg = fmt.Sprintf("У пользователя с user_id=%d нет друзей.", userIdInt)
		} else {
			successMsg = fmt.Sprintf("Список друзей пользователя с user_id=%d:\n", userIdInt)
			successMsg += strings.Join(friends, "\n")
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		log.Infof("Successfully got friends for user with user_id=%d", userIdInt)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside GetAllFriendsHandler, inappropriate http.Request.Method: GET required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
