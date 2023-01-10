package databaseRequests

import (
	"database/sql"
	"fmt"
	"strconv"
)

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
