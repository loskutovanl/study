package v1

import (
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

// ProcessInvalidRequestMethod обработка некорректного метода
func ProcessInvalidRequestMethod(w http.ResponseWriter, handlerName, methodRequired, methodReceived string) {
	log.Infof("Inside %s, inappropriate http.Request.Method: %s required, %s received", handlerName, methodRequired, methodReceived)
	w.WriteHeader(http.StatusBadRequest)
}
