package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

type Storage struct {
	repository map[int]*User
}

func NewStorage() *Storage {
	return &Storage{
		make(map[int]*User),
	}
}

func (s *Storage) GetNextId() int {
	return len(s.repository) + 1
}

func (s *Storage) MakeFriends(sourceId, targetId int) (string, error) {
	var sourceUser *User
	var targetUser *User

	for id, user := range s.repository {
		switch id {
		case sourceId:
			sourceUser = user
		case targetId:
			targetUser = user
		}
	}

	if sourceUser == nil {
		return "", fmt.Errorf("no *User found with given sourceId %d", sourceId)
	}
	if targetUser == nil {
		return "", fmt.Errorf("no *User found with given targetId %d", targetId)
	}

	s.repository[sourceId].Friends = append(s.repository[sourceId].Friends, targetUser)
	s.repository[targetId].Friends = append(s.repository[targetId].Friends, sourceUser)
	return fmt.Sprintf("%s и %s теперь друзья", sourceUser.Name, targetUser.Name), nil
}

func (s *Storage) DeleteUser(targetId int) (string, error) {
	var userToDelete *User

	for id, user := range s.repository {
		if id == targetId {
			userToDelete = user
		}
	}

	if userToDelete == nil {
		return "", fmt.Errorf("no *User found with given targetId %d", targetId)
	}

	for _, friend := range userToDelete.Friends {
		for id, user := range s.repository {
			if user == friend {
				delete(s.repository, id)
				break
			}
		}
	}
	delete(s.repository, targetId)
	return userToDelete.Name, nil
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     string  `json:"age"`
	Friends []*User `json:"friends"`
}

type Friends struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

type DeleteUser struct {
	TargetId string `json:"target_id"`
}

func main() {
	storage := NewStorage()

	r := chi.NewRouter()

	r.Post("/create", storage.CreateHandler)
	r.Post("/make_friends", storage.MakeFriendsHandler)
	r.Delete("/user", storage.DeleteHandler)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func (s *Storage) CreateHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside CreateHandler")
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside CreateHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				log.Warn("Inside CreateHandler, unable to write response after reading http.Request.Body:", err2)
				return
			}
			return
		}

		var u *User
		if err = json.Unmarshal(content, &u); err != nil {
			log.Warn("Inside CreateHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				log.Warn("Inside CreateHandler, unable to write response after unmarshalling json:", err2)
				return
			}
			return
		}

		u.ID = s.GetNextId()
		s.repository[u.ID] = u

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(strconv.Itoa(u.ID)))
		if err != nil {
			log.Warn("Inside CreateHandler, unable to write response after adding user to storage:", err)
			return
		}
		return
	}

	log.Infof("Inside CreateHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Storage) MakeFriendsHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside MakeFriendsHandler")
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				log.Warn("Inside MakeFriendsHandler, unable to write response after reading http.Request.Body:", err2)
				return
			}
			return
		}

		var f *Friends
		if err = json.Unmarshal(content, &f); err != nil {
			log.Warn("Inside MakeFriendsHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		sourceId, err := strconv.Atoi(f.SourceId)
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to convert sourceId %d to int: %s", sourceId, err)
			return
		}

		targetId, err := strconv.Atoi(f.TargetId)
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to convert targetId %d to int: %s", targetId, err)
			return
		}

		requestStatus, err := s.MakeFriends(sourceId, targetId)
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to make users %d and %d friends: %s", sourceId, targetId, err)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(requestStatus))
		if err != nil {
			log.Warnf("Inside MakeFriendsHandler, unable to write response after making users %d and %d friends: %s", sourceId, targetId, err)
		}
		return
	}

	log.Infof("Inside MakeFriendsHandler, inappropriate http.Request.Method: POST required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}

func (s *Storage) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Inside DeleteHandler")
	if r.Method == "DELETE" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Warn("Inside DeleteHandler, unable to read http.Request.Body:", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err2 := w.Write([]byte(err.Error()))
			if err2 != nil {
				log.Warn("Inside DeleteHandler, unable to write response after reading http.Request.Body:", err2)
				return
			}
			return
		}

		var du *DeleteUser
		if err = json.Unmarshal(content, &du); err != nil {
			log.Warn("Inside DeleteHandler, unable to Unmarshal json:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		targetId, err := strconv.Atoi(du.TargetId)
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable to convert targetId %d to int: %s", targetId, err)
			return
		}

		requestStatus, err := s.DeleteUser(targetId)
		fmt.Println(s.repository)
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable delete user %d: %s", targetId, err)
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(requestStatus))
		if err != nil {
			log.Warnf("Inside DeleteHandler, unable to write response afterdeleting user %d: %s", targetId, err)
		}
		return
	}

	log.Infof("Inside DeleteHandler, inappropriate http.Request.Method: DELETE required, %s received", r.Method)
	w.WriteHeader(http.StatusBadRequest)
}
