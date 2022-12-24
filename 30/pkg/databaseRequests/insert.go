package databaseRequests

import (
	"database/sql"
	"fmt"
	"strconv"

	"30/30/pkg/customStructures"
)

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

// ValidateUserAndCreateUser создает запись о пользователе в таблицу "users". Переводит возраст пользователя в числовой
// тип. Возвращает id добавленного в таблицу пользователя и ошибку.
func ValidateUserAndCreateUser(db *sql.DB, user *customStructures.User) (int, error) {
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
