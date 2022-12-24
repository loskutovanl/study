package databaseRequests

import (
	"database/sql"
	"fmt"
	"strconv"
)

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
