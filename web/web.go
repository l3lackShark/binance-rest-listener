package web

import (
	"fmt"
	"net/http"
)

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello api\n")
}

func StartServer(addr string) {
	http.HandleFunc("/api", apiHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}

}
