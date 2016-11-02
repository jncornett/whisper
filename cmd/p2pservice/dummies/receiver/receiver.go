package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/jncornett/whisper"
)

type Service struct{}

func (s *Service) Push(m *whisper.Message, reply *bool) error {
	log.Printf("received message: %+v", *m)
	*reply = true
	return nil
}

func main() {
	listenAddress := flag.String("listen", ":8082", "address to listen on")
	flag.Parse()
	svc := Service{}
	rpc.Register(&svc)
	rpc.HandleHTTP()
	ln, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(ln, nil)
}
