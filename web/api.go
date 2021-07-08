package web

import (
	"fmt"
	"log"
	"net/http"
)

func currentPriceHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprint(w, "Hello api\n")
}

func StartServer(addr string) {
	http.HandleFunc("/api", currentPriceHandler)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalln(err)
	}

}
