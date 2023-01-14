package repo

import (
	"30/internal/entity"
	"context"
)

type Repository interface {
	Migrate(ctx context.Context) error
	InsertUser(user *entity.User) (int, error)
	InsertFriends(friendId, userId int) error
	SelectUser(userId int) (entity.User, error)
	SelectFriends(sourceId, targetId int) (bool, error)
	DeleteUser(user *entity.User) error
	DeleteFriends(user *entity.User) error
	UpdateUserAge(user *entity.NewAge) error
	SelectUserFriends(user *entity.User) (friends []entity.User, err error)
}
