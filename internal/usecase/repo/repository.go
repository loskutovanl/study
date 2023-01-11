package repo

import (
	"30/internal/entity"
	"context"
)

type Repository interface {
	Migrate(ctx context.Context) error
	InsertUser(user *entity.User) (int, error)
	InsertFriends(friend string, userId int) error
	SelectUser(userId int) ([]int, error)
	SelectFriends(sourceId, targetId int) ([]int, error)
}
