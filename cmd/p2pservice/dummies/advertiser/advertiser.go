package main

import (
	"flag"
	"log"
	"net/rpc"
	"time"

	"github.com/jncornett/whisper"
)

func main() {
	transmitInterval := flag.Duration("ival", time.Second, "interval to transmit peer table")
	addr := flag.String("addr", "localhost:8082", "address to advertise")
	server := flag.String("server", "localhost:8081", "server to solicit")
	flag.Parse()
	for {
		time.Sleep(*transmitInterval)
		log.Printf("advertising peer %v to %v", *addr, *server)
		peers := []whisper.Peer{{Address: *addr, LastUpdated: time.Now()}}
		client, err := rpc.DialHTTP("tcp", *server)
		if err != nil {
			log.Print(err)
			continue
		}
		var reply bool
		err = client.Call("Service.UpdatePeers", &peers, &reply)
		if err != nil {
			log.Print(err)
		}
	}
}
