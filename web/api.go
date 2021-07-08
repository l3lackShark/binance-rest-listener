package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/l3lackShark/binance-rest-listener/database"
)

type genericHTTPError struct {
	Data string `json:"error"`
}

//This shouldn't be global in production, some in memory db has to have wrapped access to it instead. Doing for simplicity reasons
var repo database.Repository

func marshallErrorStruct(errStr string) string {
	resp := genericHTTPError{
		Data: errStr,
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		panic(err) //something very serious happened
	}
	return string(respBytes)
}

func checkRequestMethod(w *http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		errStr := fmt.Sprintf("Invalid Request type, WANT: %s, GOT: %s", method, r.Method)
		(*w).WriteHeader(http.StatusBadRequest)
		fmt.Fprint((*w), marshallErrorStruct(errStr))
		return fmt.Errorf(errStr)
	}
	return nil
}

func marshallStruct(w *http.ResponseWriter, r *http.Request, data interface{}) (string, error) {
	jsonOut, err := json.Marshal(data)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		fmt.Fprint((*w), marshallErrorStruct(fmt.Sprintf("Failed to marshall json %s", err.Error())))
		fmt.Fprint((*w), string(jsonOut))
		return "", fmt.Errorf("Failed to marshall json %s", err.Error())
	}
	return string(jsonOut), nil
}

func dayPriceHandler(w http.ResponseWriter, r *http.Request) {

	if err := checkRequestMethod(&w, r, "GET"); err != nil {
		return
	}
	q := r.URL.Query()
	date := q.Get("date")
	//check input date against our format
	if !database.DateFormat.MatchString(date) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, marshallErrorStruct("Invalid date format, want: 'DD.MM.YYYY'"))
		return
	}

	//perform a database request (this should probably be replaced with some in-memory db in production)
	doc, err := repo.FindOneByDate(database.DaatabaseName, database.CollectionName, date)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, marshallErrorStruct(fmt.Sprintf("Failed to find a document %s", err.Error())))
		return
	}

	jsonOut, err := marshallStruct(&w, r, doc)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json") //firefox dev tools
	fmt.Fprint(w, jsonOut)
}

func currentPriceHandler(w http.ResponseWriter, r *http.Request) {
	if err := checkRequestMethod(&w, r, "GET"); err != nil {
		return
	}

	//perform a database request (this will probably be replaced with some in-memory db in production)
	doc, err := repo.FindOneByDate(database.DaatabaseName, database.CollectionName, time.Now().UTC().Format("02.01.2006"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, marshallErrorStruct(fmt.Sprintf("Failed to find a document %s", err.Error())))
		return
	}
	if len(doc.Stamps) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, marshallErrorStruct(fmt.Sprintf("There is no data for today %s", err.Error())))
	}
	//get the last report
	jsonOut, err := marshallStruct(&w, r, doc.Stamps[len(doc.Stamps)-1])
	if err != nil {
		return
	}
	fmt.Fprint(w, jsonOut)
}

func StartServer(addr string) {
	var err error
	repo, err = database.New(os.Getenv("MONGO_CONN_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/api/getDay", dayPriceHandler)
	http.HandleFunc("/api/getCurrent", currentPriceHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln(err)
	}

}
