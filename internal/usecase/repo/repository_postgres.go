package repo

import (
	"30/internal/entity"
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type PostgreSQLClassicRepository struct {
	db *sql.DB
}

func NewPostgreSQLClassicRepository(db *sql.DB) *PostgreSQLClassicRepository {
	return &PostgreSQLClassicRepository{
		db: db,
	}
}

func (r *PostgreSQLClassicRepository) Migrate(ctx context.Context) error {
	query := `
	CREATE TABLE IF NOT EXISTS "friends"(
    id SERIAL PRIMARY KEY,
    user1_id INT NOT NULL,
    user2_id INT NOT NULL
	);`

	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgreSQLClassicRepository) InsertUser(user *entity.User) (int, error) {
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
	err = r.db.QueryRow(query, user.Name, age).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("unable to insert user (name %s, age %d) to database table users: %s", user.Name, age, err)
	}

	return userId, nil
}

func (r *PostgreSQLClassicRepository) InsertFriends(friend string, userId int) error {
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
