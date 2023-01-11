package http

import (
	"30/internal/databaseRequests"
	"30/internal/entity"
	"30/internal/handlers"
	"30/internal/usecase"
	"30/internal/usecase/repo"
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

func newUserRoutes(mux *chi.Mux, uc usecase.UserUseCase) {
	route := &userRoutes{uc}
	mux.Post("/create", func(w http.ResponseWriter, r *http.Request) { handlers.CreateHandler(w, r, db) })
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

		// Use case
		userUseCase := usecase.New(
			repo.NewPostgreSQLClassicRepository(ur.uc.),
		)

		userId, err := userUseCase.NewUser(u)

		//userId, err := userUsecase.NewUser(u)
		if err != nil {
			log.Errorf("Inside CreateHandler: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		log.Infof("Successfully created user (user_id %d)", userId)

		// добавление связи друзей в таблицу "friends
		for _, friend := range u.Friends {
			err = databaseRequests.MakeFriendsForCreatedUser(ur.uc.r, friend, userId)
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
