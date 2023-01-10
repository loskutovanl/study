package handlers

import (
	"30/internal/databaseRequests"
	"30/internal/entity"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

// PutUserAgeHandler обрабатывает PUT-запрос на изменение возраста пользователя. Логирует возможные ошибки.
// При успешном запросе возвращает соответствующее сообщение и статус 200.
// При неуспешном запросе возвращает статус 400 (неверный метод) ии 500 (ошибка обработки данных).
func PutUserAgeHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Info("Inside PutUserAgeHandler")

	// чтение запроса и обработка ошибок
	if r.Method == "PUT" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside PutUserAgeHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// демаршализация запроса и обработка ошибок
		var na *entity.NewAge
		if err = json.Unmarshal(content, &na); err != nil {
			log.Warn("Inside PutUserAgeHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// приведение id пользователя к числовому типу
		userIdString := chi.URLParam(r, "id")
		userIdInt, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Warnf("Inside PutUserAgeHandler, unable to convert user_id %s from string to int: %s", userIdString, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// проверка, что пользователь существует в таблице "users"
		err = databaseRequests.CheckUserExistsInUsersTable(db, userIdInt)
		if err != nil {
			log.Warnf("Inside PutUserAgeHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обновление возраста пользователя
		err = databaseRequests.UpdateUserAge(db, userIdInt, na.Age)
		if err != nil {
			log.Warnf("Inside PutUserAgeHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		successMsg := "Возраст пользователя успешно обновлён"
		log.Infof("Successfully changed user (user_id=%d) age to %s", userIdInt, na.Age)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside PutUserAgeHandler, inappropriate http.Request.Method: PUT required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
