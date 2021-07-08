package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/l3lackShark/binance-rest-listener/database"
)

//Was originally planning on using websockets, but decided to send GET requests for simplicity. Seems like they already have an endpoint for the avgPrice of the last 5 mins (/api/v3/avgPrice)

type AVGRespone struct {
	Mins  int    `json "mins"`
	Price string `json:"price"`
}

func RequestFiveMinAVG() (AVGRespone, error) {
	resp, err := http.Get("https://api.binance.com/api/v3/avgPrice?symbol=BTCUSDT")
	if err != nil {
		return AVGRespone{}, fmt.Errorf("Failed to send http req: %e", err)
	}

	if resp.StatusCode != http.StatusOK {
		return AVGRespone{}, fmt.Errorf("Received invalid response: %d", resp.StatusCode)
	}

	out := AVGRespone{}
	bytes := new(bytes.Buffer)
	if _, err := bytes.ReadFrom(resp.Body); err != nil {
		return AVGRespone{}, fmt.Errorf("Failed to send buffer to bytes: %e", err)
	}

	json.Unmarshal(bytes.Bytes(), &out)
	return out, nil
}

//Fatals are not handled, so if an error occurs, the application will exit.

func PriceLoop() {
	ticker := time.NewTicker(time.Minute * 5)
	//establish a database connection

	repo, err := database.New(os.Getenv("MONGO_CONN_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	for range ticker.C {
		resp, err := RequestFiveMinAVG()
		if err != nil {
			log.Fatalln(err)
		}
		ins := database.Document{
			Date:  time.Now().UTC().Format("02.01.2006"), //DD.MM.YYYY
			Price: resp.Price,
			Time:  time.Now().UTC().Format("15:04"), //HH:MM
		}

		err = repo.UpdateOrInsertOne(database.DaatabaseName, database.CollectionName, ins)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
