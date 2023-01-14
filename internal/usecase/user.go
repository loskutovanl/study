package usecase

import (
	"30/internal/entity"
	"30/internal/usecase/repo"
	"fmt"
	log "github.com/sirupsen/logrus"
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
	for _, friendId := range user.Friends {
		err = uc.r.InsertFriends(friendId, userId)
		if err != nil {
			log.Errorf("UserUseCase - NewUser - s.r.InsertFriends: %w", err)
		} else {
			log.Infof("Successfully added friends relation (user1_id %d, user2_id %d) to database table friends", userId, friendId)
		}
	}

	return userId, nil
}

func (uc *UserUseCase) NewFriends(friends *entity.Friends) error {
	err := uc.r.InsertFriends(friends.SourceId, friends.TargetId)
	if err != nil {
		return fmt.Errorf("UserUseCase - NewFriends - s.r.InsertFriends: %w", err)
	}

	log.Infof("Successfully added friends relation (user1_id %d, user2_id %d) to database table friends", friends.SourceId, friends.TargetId)
	return nil
}

func (uc *UserUseCase) DeleteUser(user *entity.User) (userName string, err error) {

	userFromRepo, err := uc.r.SelectUser(user.Id)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.SelectUsername: %w", err)
	}
	userName = userFromRepo.Name

	err = uc.r.DeleteUser(user)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.DeleteUser: %w", err)
	}
	log.Infof("Successfully deleted user with id = %d (name %s)", user.Id, userName)

	err = uc.r.DeleteFriends(user)
	if err != nil {
		return userName, fmt.Errorf("UserUseCase - DeleteUser - s.r.DeleteFriends: %w", err)
	}
	log.Infof("Successfully deleted friends record for user with id = %d", user.Id)

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

func (uc *UserUseCase) GetFriends(user *entity.User) (friends []entity.User, err error) {
	// проверка, что пользователь существует в таблице "users"
	_, err = uc.r.SelectUser(user.Id)
	if err != nil {
		return friends, fmt.Errorf("UserUseCase - GetFriends - s.r.SelectUser: %s", err)
	}

	// извлечение друзей пользователя из таблиц "users" и "friends"
	friends, err = uc.r.SelectUserFriends(user)
	if err != nil {
		return friends, fmt.Errorf("UserUseCase - GetFriends - s.r.SelectUserFriends: %s", err)
	}
	log.Infof("Successfully got friends for user with user_id=%d", user.Id)

	return friends, nil
}
