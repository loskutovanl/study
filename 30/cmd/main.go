package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

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

// настройка логов
func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

// User содержит информацию о пользователе: имя, возраст, список друзей
type User struct {
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

// Friends содержит информацию об id двух пользователей, отправивших запрос на дружбу
type Friends struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

// DeleteUser содержит информацию об id пользователя, на которого запрашивается удаление
type DeleteUser struct {
	TargetId string `json:"target_id"`
}

// NewAge содержит инормацию о новом возрасте пользователя
type NewAge struct {
	Age string `json:"new_age"`
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

	// создание роутера и регистрация хендлеров
	r := chi.NewRouter()
	r.Post("/create", func(w http.ResponseWriter, r *http.Request) { CreateHandler(w, r, db) })
	r.Post("/make_friends", func(w http.ResponseWriter, r *http.Request) { MakeFriendsHandler(w, r, db) })
	r.Delete("/user", func(w http.ResponseWriter, r *http.Request) { DeleteHandler(w, r, db) })
	r.Get("/friends/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { GetAllFriendsHandler(w, r, db) })
	r.Put("/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { PutUserAgeHandler(w, r, db) })

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

		// валидация данных пользователя, добавление пользователя в таблицу "users",
		userId, err := ValidateUserAndCreateUser(db, u)
		if err != nil {
			log.Errorf("Inside CreateHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully created user (user_id %d)", userId)

		// добавление связи друзей в таблицу "friends
		for _, friend := range u.Friends {
			err = MakeFriendsForCreatedUser(db, friend, userId)
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
		var f *Friends
		if err = json.Unmarshal(content, &f); err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обработка и проверка значений пользовательских id
		err = ValidateUsersIdAndMakeFriends(db, f.SourceId, f.TargetId)
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
		var du *DeleteUser
		if err = json.Unmarshal(content, &du); err != nil {
			log.Warn("Inside DeleteHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обработка пользовательского id, удаление пользователя из таблицы "users"
		deletedUserName, err := DeleteUserFromUserTable(db, du.TargetId)
		if err != nil {
			log.Errorf("Inside DeleteHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully deleted user with id = %s (name %s)", du.TargetId, deletedUserName)

		// удаление записей о друзьях удаленного пользователя
		err = DeleteFriendsFromFriendsTable(db, du.TargetId)
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
		err = CheckUserExistsInUsersTable(db, userIdInt)
		if err != nil {
			log.Warnf("Inside GetAllFriendsHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// извлечение друзей пользователя из таблиц "users" и "friends"
		friends, err := GetFriendsForUser(db, userIdInt)
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
		var na *NewAge
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
		err = CheckUserExistsInUsersTable(db, userIdInt)
		if err != nil {
			log.Warnf("Inside PutUserAgeHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// обновление возраста пользователя
		err = UpdateUserAge(db, userIdInt, na.Age)
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

// CheckUserExistsInUsersTable проверяет, содержится ли в таблице "users" пользователь с указанным userId.
// Если такого пользователя нет или возникает проблема с открытием базы данных, возвращает ошибку.
func CheckUserExistsInUsersTable(db *sql.DB, userId int) error {
	var (
		query   = `select "id" from "users" where "id" = $1`
		records []int
		record  int
	)

	rows, err := db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("unable to perform select query on users table in database: %w", err)
	}

	for rows.Next() {
		rows.Scan(&record)
		records = append(records, record)
	}
	if len(records) == 0 {
		return fmt.Errorf("unable to find users with userId %d", userId)
	}

	return nil
}

// CheckUsersAreNotFriends проверяет, не являются ли пользователи с sourceId и targetId друзьями по таблице "friends".
// Если пользователи уже являются друзьями или возникает проблема с открытием базы данных, возвращает ошибку.
func CheckUsersAreNotFriends(db *sql.DB, sourceId, targetId int) error {
	var (
		query   = `select "id" from "friends" where ("user1_id" = $1 and "user2_id" = $2) or ("user1_id" = $2 and "user2_id" = $1)`
		records []int
		record  int
	)

	rows, err := db.Query(query, sourceId, targetId)
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("unable to perform select query on friends table in database: %w", err)
	}

	for rows.Next() {
		rows.Scan(&record)
		records = append(records, record)
	}

	if len(records) != 0 {
		return fmt.Errorf("users with sourceId %d and targetId %d are already friends", sourceId, targetId)
	}

	return nil
}

// InsertUsersIntoFriendsTable делает пользователей с sourceId и targetId друзьями, делая соответствующую запись в таблицу "friends".
// Если возникает проблема с записью в базу данных, возвращает ошибку.
func InsertUsersIntoFriendsTable(db *sql.DB, sourceId, targetId int) error {
	var query = `insert into "friends"("user1_id", "user2_id") values($1, $2)`

	_, err := db.Exec(query, sourceId, targetId)
	if err != nil {
		return fmt.Errorf("unable to perform sql insert request for user1_id %d and user2_id %d into table friends", sourceId, targetId)
	}

	return nil
}

// ValidateUsersIdAndMakeFriends приводит строковые типы id пользователей sourceIdString и targetIdString к типам int.
// Проверяет, что пользователи существует в таблице "users". Проверяет, что указаны различные id пользователей.
// Проверяет, что пользователи с указанными id еще не являются друзьями в таблице "friends". Записывает пользователей
// в таблицу "friends". Если на любом этапе возникает ошибка, возвращет ее. В случае успеа возвращает nil.
func ValidateUsersIdAndMakeFriends(db *sql.DB, sourceIdString, targetIdString string) error {
	var (
		usersIdString = []string{sourceIdString, targetIdString}
		usersIdInt    = []int{0, 0}
		err           error
	)

	for i, id := range usersIdString {
		// приведение типов друзей пользователя из string в int
		usersIdInt[i], err = strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("unable to convert userId %s from string to int: %s", id, err)
		}

		// проверка, что пользователи существуют в таблице "users"
		err = CheckUserExistsInUsersTable(db, usersIdInt[i])
		if err != nil {
			return err
		}
	}

	sourceIdInt := usersIdInt[0]
	targetIdInt := usersIdInt[1]
	// проверка, что указаны различные id двух пользователей
	if sourceIdInt == targetIdInt {
		return fmt.Errorf("unable to befriend user (userId %d) with himself", sourceIdInt)
	}

	// проверка, что указанные id пользователей еще не являются друзьями
	err = CheckUsersAreNotFriends(db, sourceIdInt, targetIdInt)
	if err != nil {
		return err
	}

	// запись друзей в таблицу "friends"
	err = InsertUsersIntoFriendsTable(db, sourceIdInt, targetIdInt)
	if err != nil {
		return err
	}

	return nil
}

// ValidateUserAndCreateUser создает запись о пользователе в таблицу "users". Переводит возраст пользователя в числовой
// тип. Возвращает id добавленного в таблицу пользователя и ошибку.
func ValidateUserAndCreateUser(db *sql.DB, user *User) (int, error) {
	var (
		userId int
		query  = `insert into "users" ("name", "age") values($1, $2) returning "id"`
	)

	// приведение типов возраста пользователя, запись имени и возраста пользователя в таблицу "users"
	age, err := strconv.Atoi(user.Age)
	if err != nil {
		return userId, fmt.Errorf("unable to convert user age %s from string to int: %s", user.Age, err)
	}

	// запись пользователя в базу данных в таблицу "users"
	err = db.QueryRow(query, user.Name, age).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("unable to insert user (name %s, age %d) to database table users: %s", user.Name, age, err)
	}

	return userId, nil
}

// MakeFriendsForCreatedUser создает записи о друзьях нового пользователя с userId в таблицу "friends". Переводит
// возраст пользователя в числовое значение, проверяет что пользователь с friendId существует в таблице "users".
// Проверяет, что пользователи с id userId, friendId еще не друзья и добавляет запись о друзьях в базу данных.
// Возвращает соответствующую ошибку, если на любом этапе что-то идет не так.
func MakeFriendsForCreatedUser(db *sql.DB, friend string, userId int) error {
	var (
		query = `insert into "friends"("user1_id", "user2_id") values($1, $2)`
	)

	// перевод типа возраста пользователя в числовое значение
	friendId, err := strconv.Atoi(friend)
	if err != nil {
		return fmt.Errorf("unable to convert friendId %s to int: %s", friend, err)
	}

	// проверка, что пользователь с id friendId существует в таблице пользователей
	err = CheckUserExistsInUsersTable(db, friendId)
	if err != nil {
		return err
	}

	// проверка, что пользователи с id userId, friendId еще не друзья
	err = CheckUsersAreNotFriends(db, userId, friendId)
	if err != nil {
		return err
	}

	// добавление записи о друзьях в базу данных
	_, err = db.Exec(query, userId, friend)
	if err != nil {
		return fmt.Errorf("unable to insert friends (user1_id %d, user2_id %d) to database table friends: %s", userId, friendId, err)
	}

	return nil
}

// DeleteUserFromUserTable удаляет запис о пользователе из таблицы "users". Переводит id пользователя в числовой формат.
// Возвращает имя удаленного пользователя (полученное селект-запросом) и соответствующую ошибку, если на любом
// этапе что-то идет не так.
func DeleteUserFromUserTable(db *sql.DB, userIdString string) (string, error) {
	var (
		deletedUserName string
		queryDelete     = `delete from "users" where "id" = $1`
		querySelect     = `select distinct "name" from "users" where "id" = $1`
	)

	// перевод id пользователя в числовой формат
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		return deletedUserName, fmt.Errorf("unable to convert targetId %d to int: %s", userIdInt, err)
	}

	// получение имени удаляемого пользователя с помощью селект-запроса
	rows, err := db.Query(querySelect, userIdInt)
	if err != nil {
		return deletedUserName, fmt.Errorf("unable to get user name for user_id = %d", userIdInt)
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&deletedUserName)
		if err != nil {
			return deletedUserName, err
		}
	}

	// удаление пользователя из таблицы "users"
	_, err = db.Exec(queryDelete, userIdInt)
	if err != nil {
		return deletedUserName, err
	}

	return deletedUserName, nil
}

// DeleteFriendsFromFriendsTable удаляет записи о всех друзьях пользователя с userIdString из таблицы "friends".
// Возвращает ошибку, если что-то идет не так.
func DeleteFriendsFromFriendsTable(db *sql.DB, userIdString string) error {
	var (
		query = `delete from "friends" where "user1_id" = $1 or "user2_id" = $1`
	)

	// перевод id пользователя в числовой формат
	userIdInt, err := strconv.Atoi(userIdString)
	if err != nil {
		return fmt.Errorf("unable to convert targetId %d to int: %s", userIdInt, err)
	}

	// удаление записей о друзьях из таблицы "friends"
	_, err = db.Exec(query, userIdInt)
	if err != nil {
		return fmt.Errorf("unable to delete from friends where user1_id or user2_id equal to %s", userIdString)
	}

	return nil
}

// GetFriendsForUser извлекает записи о всех друзьях пользователя с userIdInt из объединения таблиц "friends" и "users".
// Возвращает ошибку при неуспешном запросе.
func GetFriendsForUser(db *sql.DB, userIdInt int) ([]string, error) {
	var (
		query = `select "name", "age" from "users" inner join "friends" on users.id = friends.user1_id where user2_id = $1 union select "name", "age" from "users" inner join "friends" on users.id = friends.user2_id where user1_id = $1`
		users []string
		name  string
		age   int
	)

	rows, err := db.Query(query, userIdInt)
	defer rows.Close()
	if err != nil {
		return users, fmt.Errorf("unable to perform select query on getting friends for user_id %d: %s", userIdInt, err)
	}

	for rows.Next() {
		rows.Scan(&name, &age)
		users = append(users, strings.Join([]string{name, strconv.Itoa(age)}, " "))
	}

	return users, nil
}

// UpdateUserAge меняет возраст пользователя с userId на новое значение newAge.
// Возвращает ошибку при неуспешном запросе.
func UpdateUserAge(db *sql.DB, userIdInt int, newAgeString string) error {
	query := `update "users" set "age" = $1 where "id" = $2`

	newAgeInt, err := strconv.Atoi(newAgeString)
	if err != nil {
		return fmt.Errorf("unable to convert age %s from string to int: %s", newAgeString, err)
	}

	_, err = db.Exec(query, newAgeInt, userIdInt)
	if err != nil {
		return fmt.Errorf("unable to update age of user with suer_id=%d: %s", userIdInt, err)
	}

	return nil
}
