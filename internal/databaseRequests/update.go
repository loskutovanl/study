package databaseRequests

import (
	"database/sql"
	"fmt"
	"strconv"
)

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
