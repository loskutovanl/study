package usecase

import (
	"30/internal/entity"
	"30/internal/usecase/repo"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
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
	// добавление нового пользователя в таблицу "users"
	userId, err := uc.r.InsertUser(user)
	if err != nil {
		return userId, fmt.Errorf("UserUseCase - NewUser - s.r.InsertUser: %w", err)
	}
	log.Infof("Successfully created user (user_id %d)", userId)

	// добавление связи друзей в таблицу "friends"
	for _, friend := range user.Friends {
		friendId, err := strconv.Atoi(friend)
		if err != nil {
			log.Errorf("UserUseCase - NewUser - s.r.InsertFriends:unable to convert friendId %s to int: %s", friend, err)
		}
		err = uc.r.InsertFriends(friendId, userId)
		if err != nil {
			log.Errorf("UserUseCase - NewUser - s.r.InsertFriends: %s", err)
		} else {
			log.Infof("Successfully added friends relation (user1_id %d, user2_id %s) to database table friends", userId, friend)
		}
	}

	return userId, nil
}

func (uc *UserUseCase) NewFriends(friends *entity.Friends) error {
	err := uc.r.InsertFriends(friends.SourceId, friends.TargetId)
	if err != nil {
		return fmt.Errorf("UserUseCase - NewFriends - s.r.InsertFriends: %s", err)
	}

	log.Infof("Successfully added friends relation (user1_id %d, user2_id %d) to database table friends", friends.SourceId, friends.TargetId)
	return nil
}

func (uc *UserUseCase) DeleteUser(user *entity.DeleteUser) (userName string, err error) {

	userName, err = uc.r.SelectUsername(user)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.SelectUsername: %s", err)
	}

	err = uc.r.DeleteUser(user)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.DeleteUser: %s", err)
	}
	log.Infof("Successfully deleted user with id = %d (name %s)", user.TargetId, userName)

	err = uc.r.DeleteFriends(user)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.DeleteFriends: %s", err)
	}
	log.Infof("Successfully deleted friends record for user with id = %d", user.TargetId)

	return userName, nil
}

func (uc *UserUseCase) UpdateUserAge(user *entity.NewAge) error {
	// проверка, что пользователь существует в таблице "users"
	_, err := uc.r.SelectUser(user.Id)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdateUserAge - s.r.SelectUser: %s", err)
	}

	// обновление возраста пользователя
	err = uc.r.UpdateUserAge(user)
	if err != nil {
		return fmt.Errorf("UserUseCase - UpdateUserAge - s.r.UpdateUserAge: %s", err)
	}
	log.Infof("Successfully changed user (user_id=%d) age to %d", user.Id, user.Age)

	return nil
}
