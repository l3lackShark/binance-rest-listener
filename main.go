package main

import (
	"os"

	"github.com/l3lackShark/binance-rest-listener/envvars"
	"github.com/l3lackShark/binance-rest-listener/web"
)

var (
	ServerAddr string = "localhost:24080"
)

func init() {
	envvars.LoadEnv()
	if os.Getenv("SERVER_IP") != "" {
		ServerAddr = os.Getenv("SERVER_IP")
	}
}

func main() {
	go web.PriceLoop()
	web.StartServer(ServerAddr)
}
