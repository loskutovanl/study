package v1

import (
	"30/internal/entity"
	"30/internal/usecase"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type userRoutes struct {
	uc usecase.UserUseCase
}

func NewUserRoutes(mux *chi.Mux, uc *usecase.UserUseCase) {
	ur := &userRoutes{*uc}
	mux.Post("/create", func(w http.ResponseWriter, r *http.Request) { ur.createUser(w, r) })
}

type userRequest struct {
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

func (ur *userRoutes) createUser(w http.ResponseWriter, r *http.Request) {
	var handlerName = "createUser"
	log.Infof("Inside %s", handlerName)

	if r.Method == "POST" {
		content, err := ReadHttpRequest(w, r, handlerName)
		if err != nil {
			return
		}

		// UnmarshalHttpRequest демаршализация запроса и обработка ошибок
		var userRequest *userRequest
		if err = json.Unmarshal(content, &userRequest); err != nil {
			log.Warn("Inside %s, unable to Unmarshal json: %s", handlerName, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// приведение типов возраста пользователя, запись имени и возраста пользователя в таблицу "users"
		ageInt, err := strconv.Atoi(userRequest.Age)
		if err != nil {
			log.Errorf("unable to convert user age %s from string to int: %s", userRequest.Age, err)
			return
		}

		// добавление пользователя в таблицу "users",
		userId, err := ur.uc.NewUser(&entity.User{
			Name:    userRequest.Name,
			Age:     ageInt,
			Friends: userRequest.Friends,
		})

		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully created user (user_id %d)", userId)

		// добавление связи друзей в таблицу "friends
		for _, friend := range user.Friends {
			err = ur.uc.NewUserFriends(friend, userId)
			if err != nil {
				log.Errorf("Inside CreateHandler: %s", err)
			} else {
				log.Infof("Successfully added friends relation (user1_id %d, user2_id %s) to database table friends", userId, friend)
			}
		}

		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(strconv.Itoa(userId)))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, "POST", r.Method)
}
