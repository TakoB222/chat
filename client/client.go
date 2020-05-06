package main

import (
	"bufio"
	"flag"
	"fmt"
	"golang.org/x/net/websocket"
	"math/rand"
	"os"
	"time"
)

type Message struct {
	Text string `json:"text"`
}

var (
	port = flag.String("port", "8080", "port that using for connection")
)

func main() {
	fmt.Println("...")
	flag.Parse()

	ws, err := connect()
	if err != nil {
		fmt.Println("connection problem")
		return
	}
	var m Message

	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())
				break
			}

			fmt.Println("Messsage: ", m)
		}

	}()

	scaner := bufio.NewScanner(os.Stdin)

	for scaner.Scan() {
		text := scaner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: text,
		}

		err := websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Error sending message: ", err.Error())
			break
		}

	}
}

func connect() (ws *websocket.Conn, err error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", *port), "", makeIP())
}

func makeIP() string {
	var ip [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		ip[i] = rand.Intn(256)
	}

	return fmt.Sprintf("http://%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3])
}
