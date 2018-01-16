package main

import (
	"flag"
	"log"

	"github.com/jchv/eventgun/evgun"
)

var from, to, address string

func init() {
	flag.StringVar(&address, "address", ":4547", "Address to listen on.")
	flag.Parse()

	from = flag.Arg(0)
	to = flag.Arg(1)

	if from == "" || to == "" {
		log.Fatalln("Please enter arguments for 'from' and 'to'.")
	}
}

func main() {
	server := evgun.NewNotifyServer(from, to)
	log.Fatalln(server.Listen(address))
}
