package v1

import (
	"30/internal/entity"
	"30/internal/usecase"
	"encoding/json"
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
	mux.Post("/users/new", func(w http.ResponseWriter, r *http.Request) { ur.createUser(w, r) })
	mux.Post("/users/befriend", func(w http.ResponseWriter, r *http.Request) { ur.makeFriends(w, r) })
	mux.Delete("/users/delete", func(w http.ResponseWriter, r *http.Request) { ur.deleteUser(w, r) })
	mux.Get("/users/{id:[0-9]+}/friends", func(w http.ResponseWriter, r *http.Request) { ur.getFriends(w, r) })
	mux.Put("/users/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) { ur.updateUserAge(w, r) })
}

type userRequest struct {
	Name    string   `json:"name"`
	Age     string   `json:"age"`
	Friends []string `json:"friends"`
}

type userResponse struct {
	Id int `json:"id"`
}

func (ur *userRoutes) createUser(w http.ResponseWriter, r *http.Request) {
	var (
		handlerName    = "createUser"
		methodRequired = "POST"
	)
	log.Infof("Inside %s", handlerName)

	if r.Method == methodRequired {
		content, err := ReadHttpRequest(w, r, handlerName)
		if err != nil {
			return
		}

		var request *userRequest
		err = UnmarshalRequest(w, content, handlerName, &request)
		if err != nil {
			return
		}

		// приведение типов возраста пользователя и списка друзей, запись имени и возраста пользователя в таблицу "users"
		ageInt, err := strconv.Atoi(request.Age)
		if err != nil {
			log.Errorf("unable to convert user age %s from string to int: %s", request.Age, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		var (
			friendsArrayInt []int
			friendInt       int
		)
		for _, friendString := range request.Friends {
			friendInt, err = strconv.Atoi(friendString)
			if err != nil {
				log.Errorf("unable to convert friends (friend_id %s) from string to int: %s", friendString, err)
				friendsArrayInt = nil
				break
			}
			friendsArrayInt = append(friendsArrayInt, friendInt)
		}

		// добавление пользователя в таблицу "users"
		userId, err := ur.uc.NewUser(&entity.User{
			Name:    request.Name,
			Age:     ageInt,
			Friends: friendsArrayInt,
		})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		data := userResponse{userId}
		err = json.NewEncoder(w).Encode(data)
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, methodRequired, r.Method)
}

type friendsRequest struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

func (ur *userRoutes) makeFriends(w http.ResponseWriter, r *http.Request) {
	var (
		handlerName    = "makeFriends"
		methodRequired = "POST"
	)
	log.Infof("Inside %s", handlerName)

	if r.Method == methodRequired {
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
			ProcessStatusInternalServerError(w, err)
			return
		}
		targetId, err := strconv.Atoi(request.TargetId)
		if err != nil {
			log.Errorf("unable to convert targetId %s from string to int: %s", request.TargetId, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		err = ur.uc.NewFriends(&entity.Friends{
			SourceId: sourceId,
			TargetId: targetId,
		})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		successMsg := fmt.Sprintf("%d и %d теперь друзья", sourceId, targetId)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, methodRequired, r.Method)
}

type deleteUserRequest struct {
	TargetId string `json:"target_id"`
}

func (ur *userRoutes) deleteUser(w http.ResponseWriter, r *http.Request) {
	var (
		handlerName    = "deleteUser"
		methodRequired = "DELETE"
	)
	log.Infof("Inside %s", handlerName)

	if r.Method == methodRequired {
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
			ProcessStatusInternalServerError(w, err)
			return
		}

		// обработка пользовательского id, удаление пользователя из таблицы "users" & "friends"
		deletedUserName, err := ur.uc.DeleteUser(&entity.User{Id: targetId})
		if err != nil {
			log.Errorf("Inside %s: %s", handlerName, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(deletedUserName))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, methodRequired, r.Method)
}

type updateAgeRequest struct {
	Age string `json:"new_age"`
}

func (ur *userRoutes) updateUserAge(w http.ResponseWriter, r *http.Request) {
	var (
		handlerName    = "updateUserAge"
		methodRequired = "PUT"
	)
	log.Infof("Inside %s", handlerName)

	if r.Method == methodRequired {
		content, err := ReadHttpRequest(w, r, handlerName)
		if err != nil {
			return
		}

		var request *updateAgeRequest
		err = UnmarshalRequest(w, content, handlerName, &request)
		if err != nil {
			return
		}

		// приведение id пользователя, возраста к числовому типу
		userIdString := chi.URLParam(r, "id")
		userIdInt, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Warnf("Inside %s, unable to convert user_id %s from string to int: %s", handlerName, userIdString, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		ageInt, err := strconv.Atoi(request.Age)
		if err != nil {
			log.Warnf("unable to convert age %s from string to int: %s", request.Age, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		err = ur.uc.UpdateUserAge(&entity.NewAge{
			Id:  userIdInt,
			Age: ageInt,
		})
		if err != nil {
			ProcessStatusInternalServerError(w, err)
			return
		}

		// вывод сообщения об успехе в случае отсутствия ошибок
		successMsg := "Возраст пользователя успешно обновлён"
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(successMsg))
		return
	}

	ProcessInvalidRequestMethod(w, handlerName, methodRequired, r.Method)
}

type friendResponse struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type friendsResponse struct {
	Friend []friendResponse
}

func (ur *userRoutes) getFriends(w http.ResponseWriter, r *http.Request) {
	var (
		handlerName    = "GetFriends"
		methodRequired = "GET"
	)
	log.Infof("Inside %s", handlerName)

	if r.Method == methodRequired {
		// приведение userId к числовому типу и обработка ошибок
		userIdString := chi.URLParam(r, "id")
		userIdInt, err := strconv.Atoi(userIdString)
		if err != nil {
			log.Warnf("Inside %s, unable to convert user_id %s from string to int: %s", handlerName, userIdString, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		var friends []entity.User
		friends, err = ur.uc.GetFriends(&entity.User{
			Id: userIdInt,
		})
		if err != nil {
			log.Warnf("Inside %s: %s", handlerName, err)
			ProcessStatusInternalServerError(w, err)
			return
		}

		// запись ответа и вывод
		data := friendsResponse{}
		if len(friends) != 0 {
			for _, friend := range friends {
				data.Friend = append(data.Friend, friendResponse{
					Id:   friend.Id,
					Age:  friend.Age,
					Name: friend.Name,
				})
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(data)
		return

	}

	ProcessInvalidRequestMethod(w, handlerName, methodRequired, r.Method)
}
