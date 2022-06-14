package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

func InitBackendApi() {
	http.HandleFunc("/API/get_post", apiGetPost)

}

func ParseName(w http.ResponseWriter, r *http.Request) (string, bool) {
	keys, ok := r.URL.Query()["name"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "arguments params are missing"}`))
		return " ", false
	}
	postName := keys[0]

	return postName, true
}

func SendData(data interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(data)
	if err != nil {
		logrus.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't marshal json"}`))
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func isMethodGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return false
	}
	return true
}

//func main() {
//	mydb.ConnectToDb()
//	logrus.Info("Connected to db")
//	InitBackendApi()
//	logrus.Info("API initialised")
//
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		mydb.UppendErrorWithPath(err)
//		logrus.Info("Listen And Serve error", err)
//		return
//	}
//}
