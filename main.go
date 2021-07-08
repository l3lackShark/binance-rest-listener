package main

import (
	"flag"

	"github.com/l3lackShark/binance-ws-listener/envvars"
	"github.com/l3lackShark/binance-ws-listener/web"
)

var (
	ServerAddr = flag.String("serverip", "localhost:24080", "http service address")
)

func init() {
	envvars.LoadEnv()
	flag.Parse()
}

func main() {
	go web.PriceLoop()
	web.StartServer(*ServerAddr)

}
