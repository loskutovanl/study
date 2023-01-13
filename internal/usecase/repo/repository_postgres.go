package repo

import (
	"30/internal/entity"
	"context"
	"database/sql"
	"fmt"
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

	err := r.db.QueryRow(query, user.Name, user.Age).Scan(&userId)
	if err != nil {
		return userId, fmt.Errorf("unable to insert user (name %s, age %d) to database table users: %s", user.Name, user.Age, err)
	}

	return userId, nil
}

func (r *PostgreSQLClassicRepository) InsertFriends(friendId, userId int) error {
	var (
		query = `insert into "friends"("user1_id", "user2_id") values($1, $2)`
	)

	// проверка, что пользователь с id friendId существует в таблице пользователей
	_, err := r.SelectUser(friendId)
	if err != nil {
		return err
	}

	// проверка, что пользователи с id userId, friendId еще не друзья
	_, err = r.SelectFriends(userId, friendId)
	if err != nil {
		return err
	}

	// добавление записи о друзьях в базу данных
	_, err = r.db.Exec(query, userId, friendId)
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

func (r *PostgreSQLClassicRepository) DeleteUser(user *entity.DeleteUser) error {
	var queryDelete = `delete from "users" where "id" = $1`

	_, err := r.db.Exec(queryDelete, user.TargetId)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgreSQLClassicRepository) DeleteFriends(user *entity.DeleteUser) error {
	var query = `delete from "friends" where "user1_id" = $1 or "user2_id" = $1`

	_, err := r.db.Exec(query, user.TargetId)
	if err != nil {
		return fmt.Errorf("unable to delete from friends where user1_id or user2_id equal to %s", user.TargetId)
	}

	return nil
}

func (r *PostgreSQLClassicRepository) SelectUsername(user *entity.DeleteUser) (userName string, err error) {
	var querySelect = `select distinct "name" from "users" where "id" = $1`

	rows, err := r.db.Query(querySelect, user.TargetId)
	if err != nil {
		return userName, fmt.Errorf("unable to get user name for user_id = %d", user.TargetId)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&userName)
		if err != nil {
			return userName, err
		}
	}

	return userName, nil
}

func (r *PostgreSQLClassicRepository) UpdateUserAge(user *entity.NewAge) error {
	var query = `update "users" set "age" = $1 where "id" = $2`

	_, err := r.db.Exec(query, user.Age, user.Id)
	if err != nil {
		return fmt.Errorf("unable to update age of user with user_id=%d: %s", user.Id, err)
	}

	return nil
}
