package repo

import (
	"30/internal/entity"
	"context"
)

type Repository interface {
	Migrate(ctx context.Context) error
	InsertUser(user *entity.User) (int, error)
	InsertFriends(friendId, userId int) error
	SelectUser(userId int) ([]int, error)
	SelectFriends(sourceId, targetId int) ([]int, error)
	DeleteUser(user *entity.User) error
	DeleteFriends(user *entity.User) error
	SelectUsername(user *entity.User) (userName string, err error)
	UpdateUserAge(user *entity.NewAge) error
	SelectUserFriends(user *entity.User) (friends []entity.User, err error)
}
