package main

import (
	"flag"
	"log"
	"math/rand"
	"net/rpc"
	"time"

	"github.com/jncornett/whisper"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randMessage() whisper.Message {
	return whisper.Message{
		Author:  randString(10),
		Token:   randString(24),
		Content: randString(140),
		TTL:     uint(rand.Uint32()),
	}
}

func main() {
	addr := flag.String("addr", "localhost:8081", "address of the p2p service")
	ival := flag.Duration("ival", time.Second, "interval of messages")
	for {
		log.Print("starting loop")
		time.Sleep(*ival)
		m := randMessage()
		log.Printf("dialing %v", *addr)
		client, err := rpc.DialHTTP("tcp", *addr)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Printf("calling Service.Push(%+v)", &m)
		var reply bool
		err = client.Call("Service.Push", &m, &reply)
		if err != nil {
			log.Print(err)
		}
		log.Printf("got reply %v", reply)
	}
}
