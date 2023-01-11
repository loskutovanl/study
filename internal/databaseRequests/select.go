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
