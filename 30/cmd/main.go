package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Service struct {
	storage map[int]*User
}

func NewService() *Service {
	return &Service{
		make(map[int]*User),
	}
}

func (s *Service) getNextId() int {
	return len(s.storage) + 1
}

func (s *Service) getUserById() *User {
	# TODO
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     string     `json:"age"`
	Friends []*User `json:"friends"`
}

type Friends struct {
	SourceId string `json:"source_id"`
	TargetId string `json:"target_id"`
}

func main() {
	srv := NewService()

	r := chi.NewRouter()
	//r.Use(middleware.Logger)

	r.Post("/create", srv.createHandler)
	r.Post("/make_friends", srv.makeFriendsHandler)

	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func (s *Service) createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		var u *User
		if err = json.Unmarshal(content, &u); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		u.ID = s.getNextId()
		s.storage[u.ID] = u

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte(strconv.Itoa(u.ID)))
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}


func (s *Service) makeFriendsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		var f *Friends
		if err = json.Unmarshal(content, &f); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		sourceId, _ := strconv.Atoi(f.SourceId)
		targetId, _ := strconv.Atoi(f.TargetId)
		userSource :=

	}

	w.WriteHeader(http.StatusBadRequest)
}
