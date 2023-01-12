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

	// запись пользователя в базу данных в таблицу "users"
	err := r.db.QueryRow(query, user.Name, user.Age).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("unable to insert user (name %s, age %d) to database table users: %s", user.Name, user.Age, err)
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
	_, err = r.SelectUser(friendId)
	if err != nil {
		return err
	}

	// проверка, что пользователи с id userId, friendId еще не друзья
	_, err = r.SelectFriends(userId, friendId)
	if err != nil {
		return err
	}

	// добавление записи о друзьях в базу данных
	_, err = r.db.Exec(query, userId, friend)
	if err != nil {
		return fmt.Errorf("unable to insert friends (user1_id %d, user2_id %d) to database table friends: %s", userId, friendId, err)
	}

	return nil
}

func (r *PostgreSQLClassicRepository) SelectUser(userId int) ([]int, error) {
	var (
		query   = `select "id" from "users" where "id" = $1`
		records []int
		record  int
	)

	rows, err := r.db.Query(query, userId)
	defer rows.Close()
	if err != nil {
		return records, fmt.Errorf("unable to perform select query on users table in database: %w", err)
	}

	for rows.Next() {
		rows.Scan(&record)
		records = append(records, record)
	}
	if len(records) == 0 {
		return records, fmt.Errorf("unable to find users with userId %d", userId)
	}

	return records, nil
}

func (r *PostgreSQLClassicRepository) SelectFriends(sourceId, targetId int) ([]int, error) {
	var (
		query = `select "id" from "friends" 
            	where ("user1_id" = $1 and "user2_id" = $2) 
        		or ("user1_id" = $2 and "user2_id" = $1)`
		records []int
		record  int
	)

	rows, err := r.db.Query(query, sourceId, targetId)
	defer rows.Close()
	if err != nil {
		return records, fmt.Errorf("unable to perform select query on friends table in database: %w", err)
	}

	for rows.Next() {
		rows.Scan(&record)
		records = append(records, record)
	}

	if len(records) != 0 {
		return records, fmt.Errorf("users with sourceId %d and targetId %d are already friends", sourceId, targetId)
	}

	return records, nil
}
