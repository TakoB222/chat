package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Message struct {
	Text string `json:"text"`
}

type hub struct {
	clients map[string]*websocket.Conn
	addClientChan chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan chan Message
}

var(
	port = flag.String("port", "8080", "port used for ws connection")
)

func main(){
	flag.Parse()
	log.Fatal(server(*port))
}

func server(port string)error{
	fmt.Println("server is raning...")
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))
	m := &autocert.Manager{
		Cache:      autocert.DirCache("golang-autocert"),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("ws://localhost"),
	}
	server := http.Server{
		Addr:":"+port,
		Handler:mux,
		TLSConfig:m.TLSConfig(),
	}

	return server.ListenAndServeTLS("","")
}

func handler(ws *websocket.Conn, h* hub){
	go h.run()

	h.addClientChan <- ws

	for{
		var m Message

		err := websocket.JSON.Receive(ws, &m)
		if err != nil{
			h.broadcastChan <- Message{err.Error()}
			h.removeClientChan <- ws
			return
		}

		h.broadcastChan <- m
	}
}


func newHub()*hub{
	return &hub{
		clients:make(map[string]*websocket.Conn),
		addClientChan: make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan: make(chan Message),
	}
}

func (h *hub)run(){
	for{
		select{
			case conn := <-h.addClientChan:
				h.addClients(conn)
			case conn := <-h.removeClientChan:
				h.removeClients(conn)
			case m := <-h.broadcastChan:
				h.broadcast(&m)

		}}
}

func (h *hub)addClients(conn *websocket.Conn){
	h.clients[conn.RemoteAddr().String()] = conn
}

func (h *hub)removeClients(conn *websocket.Conn){
	delete(h.clients, conn.LocalAddr().String())
}

func (h *hub)broadcast(m *Message){
	for _, conn := range h.clients{
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}

}
























