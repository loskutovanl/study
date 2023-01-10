package handlers

import (
	"30/internal/databaseRequests"
	"30/internal/entity"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// DeleteHandler обрабатывает DELETE-запрос на удаление пользователя из таблицы "users", а также стирает записи о его
// друзьях из таблицы "friends". Логирует возможные ошибки. При успешном запросе возвращает имя удаленного пользователя
// и статус 200. При неуспешном запросе возвращает статус 400 (неверный метод) ии 500 (ошибка обработки данных).
func DeleteHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	log.Info("Inside DeleteHandler")

	// чтение запроса и обработка ошибок
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside DeleteHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// демаршализация запроса и обработка ошибок
		var du *entity.DeleteUser
		if err = json.Unmarshal(content, &du); err != nil {
			log.Warn("Inside DeleteHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обработка пользовательского id, удаление пользователя из таблицы "users"
		deletedUserName, err := databaseRequests.DeleteUserFromUserTable(db, du.TargetId)
		if err != nil {
			log.Errorf("Inside DeleteHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully deleted user with id = %s (name %s)", du.TargetId, deletedUserName)

		// удаление записей о друзьях удаленного пользователя
		err = databaseRequests.DeleteFriendsFromFriendsTable(db, du.TargetId)
		if err != nil {
			log.Errorf("Inside DeleteHandler: %s", err)
		} else {
			log.Infof("Successfully deleted friends record for user with id = %s", du.TargetId)
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(deletedUserName))
		return
	}

	// обработка некорректного метода
	log.Infof("Inside DeleteHandler, inappropriate http.Request.Method: DELETE required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
