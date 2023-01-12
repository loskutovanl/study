package v1

import (
	"30/internal/entity"
	"30/internal/usecase"
	"fmt"
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
	mux.Post("/make_friends", func(w http.ResponseWriter, r *http.Request) { ur.makeFriends(w, r) })
	mux.Delete("/user", func(w http.ResponseWriter, r *http.Request) { ur.deleteUser(w, r) })
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

		var request *userRequest
		err = UnmarshalRequest(w, content, handlerName, &request)
		if err != nil {
			return
		}

		// приведение типов возраста пользователя, запись имени и возраста пользователя в таблицу "users"
		ageInt, err := strconv.Atoi(request.Age)
		if err != nil {
			log.Errorf("unable to convert user age %s from string to int: %s", request.Age, err)
			return
		}

		// добавление пользователя в таблицу "users",
		userId, err := ur.uc.NewUser(&entity.User{
			Name:    request.Name,
			Age:     ageInt,
			Friends: request.Friends,
		})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(strconv.Itoa(userId)))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, "POST", r.Method)
}

type friendsRequest struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

func (ur *userRoutes) makeFriends(w http.ResponseWriter, r *http.Request) {
	var handlerName = "makeFriends"
	log.Infof("Inside %s", handlerName)

	if r.Method == "POST" {
		content, err := ReadHttpRequest(w, r, handlerName)
		if err != nil {
			return
		}

		var request *friendsRequest
		err = UnmarshalRequest(w, content, handlerName, &request)
		if err != nil {
			return
		}

		// приведение типов id пользователя
		sourceId, err := strconv.Atoi(request.SourceId)
		if err != nil {
			log.Errorf("unable to convert sourceId %s from string to int: %s", request.SourceId, err)
			return
		}
		targetId, err := strconv.Atoi(request.TargetId)
		if err != nil {
			log.Errorf("unable to convert targetId %s from string to int: %s", request.TargetId, err)
			return
		}

		err = ur.uc.NewFriends(&entity.Friends{
			SourceId: sourceId,
			TargetId: targetId,
		})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		successMsg := fmt.Sprintf("%d и %d теперь друзья", sourceId, targetId)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, "POST", r.Method)
}

type deleteUserRequest struct {
	TargetId string `json:"target_id"`
}

func (ur *userRoutes) deleteUser(w http.ResponseWriter, r *http.Request) {
	var handlerName = "deleteUser"
	log.Infof("Inside %s", handlerName)

	if r.Method == "DELETE" {
		content, err := ReadHttpRequest(w, r, handlerName)
		if err != nil {
			return
		}

		var request *deleteUserRequest
		err = UnmarshalRequest(w, content, handlerName, &request)
		if err != nil {
			return
		}

		// приведение типа id пользователя
		targetId, err := strconv.Atoi(request.TargetId)
		if err != nil {
			log.Errorf("unable to convert user targetId %s from string to int: %s", request.TargetId, err)
			return
		}

		// обработка пользовательского id, удаление пользователя из таблицы "users" & "friends"
		deletedUserName, err := ur.uc.DeleteUser(&entity.DeleteUser{TargetId: targetId})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(deletedUserName))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, "DELETE", r.Method)
}
