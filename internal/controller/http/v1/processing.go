package v1

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// ReadHttpRequest чтение запроса и обработка ошибок
func ReadHttpRequest(w http.ResponseWriter, r *http.Request, handlerName string) ([]byte, error) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warnf("Inside %s, unable to read http.Request.Body: %s", handlerName, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return content, err
	}
	return content, nil
}

// UnmarshalRequest демаршализация запроса и обработка ошибок
func UnmarshalRequest(w http.ResponseWriter, content []byte, handlerName string, request interface{}) error {
	if err := json.Unmarshal(content, &request); err != nil {
		log.Warn("Inside %s, unable to Unmarshal json: %s", handlerName, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return err
	}
	return nil
}

// ProcessInvalidRequestMethod обработка некорректного метода
func ProcessInvalidRequestMethod(w http.ResponseWriter, handlerName, methodRequired, methodReceived string) {
	log.Infof("Inside %s, inappropriate http.Request.Method: %s required, %s received", handlerName, methodRequired, methodReceived)
	w.WriteHeader(http.StatusBadRequest)
}
