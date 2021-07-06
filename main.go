package main

import (
	"flag"

	"github.com/l3lackShark/binance-ws-listener/web"
)

var (
	ServerAddr = flag.String("serverip", "localhost:24080", "http service address")
)

func init() {
	flag.Parse()
}

func main() {
	web.StartServer(*ServerAddr)
}
