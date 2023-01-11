package v1

import (
	"30/internal/entity"
	"30/internal/usecase"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
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

func (ur *userRoutes) createUser(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside CreateHandler")

	if r.Method == "POST" {
		// чтение запроса и обработка ошибок
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside CreateHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// демаршализация запроса и обработка ошибок
		var u *entity.User
		if err = json.Unmarshal(content, &u); err != nil {
			log.Warn("Inside CreateHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// валидация данных пользователя, добавление пользователя в таблицу "users",
		userId, err := ur.uc.NewUser(u)

		if err != nil {
			log.Errorf("Inside CreateHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully created user (user_id %d)", userId)

		// добавление связи друзей в таблицу "friends
		for _, friend := range u.Friends {
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

	// обработка некорректного метода
	log.Infof("Inside CreateHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
