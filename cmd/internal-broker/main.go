package main

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"log"
	"os"
	"time"
)

const (
	EvNewMessage = iota
	EvNewAccount
)

func client() {
	sock, err := zmq.NewSocket(zmq.PULL)
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()

	if err := sock.Connect("ipc:///Users/reymons/Desktop/test.ipc"); err != nil {
		log.Fatal(err)
	}

	sock.SetSubscribe("")

	for {
		time.Sleep(2 * time.Second)
		msg, err := sock.Recv(0)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("MESSAGE: %s\n", msg)
	}
}

func server() {
	sock, err := zmq.NewSocket(zmq.PUSH)
	if err != nil {
		log.Fatal(err)
	}
	defer sock.Close()

	if err := sock.Bind("ipc:///Users/reymons/Desktop/test.ipc"); err != nil {
		log.Fatal(err)
	}

	var i int
	time.Sleep(2 * time.Second)
	for {
		if _, err := sock.Send(fmt.Sprintf("Hello %d", i), 0); err != nil {
			log.Fatal(err)
		}
		log.Printf(fmt.Sprintf("Sent %d\n", i))
		i += 1
		if i >= 5 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "client" {
		go client()
	} else {
		go server()
	}

	time.Sleep(time.Hour)
}
