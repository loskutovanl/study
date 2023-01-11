package usecase

import (
	"30/internal/entity"
	"30/internal/usecase/repo"
	"fmt"
)

type UserUseCase struct {
	r repo.Repository
}

func New(r repo.Repository) *UserUseCase {
	return &UserUseCase{
		r: r,
	}
}

func (uc *UserUseCase) NewUser(user *entity.User) (int, error) {
	userId, err := uc.r.InsertUser(user)
	if err != nil {
		return userId, fmt.Errorf("UserUseCase - NewUser - s.r.InsertUser: %w", err)
	}
	return userId, nil
}

func (uc *UserUseCase) NewUserFriends(friend string, userId int) error {
	err := uc.r.InsertFriends(friend, userId)
	if err != nil {
		return fmt.Errorf("UserUseCase - NewUserFriends - s.r.InsertFriends: %w", err)
	}
	return nil
}
