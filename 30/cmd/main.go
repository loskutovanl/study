package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type service struct {
	storage map[int]*User
}

func (s *service) getNextId() int {
	return len(s.storage) + 1
}

type User struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Friends []*User `json:"friends"`
}

func main() {
	mux := http.NewServeMux()
	srv := service{make(map[int]*User)}
	mux.HandleFunc("/create", srv.Create)

	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func (s *service) Create(rw http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, err := rw.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		var u *User
		if err = json.Unmarshal(content, &u); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			_, err := rw.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}
		u.ID = s.getNextId()
		s.storage[u.ID] = u

		rw.WriteHeader(http.StatusCreated)
		_, err = rw.Write([]byte(strconv.Itoa(u.ID)))
		if err != nil {
			return
		}
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
}
