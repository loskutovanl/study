package databaseRequests

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

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
