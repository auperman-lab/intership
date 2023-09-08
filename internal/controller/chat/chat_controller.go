package chat

import (
	"crypto/rand"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func randomUint() uint {
	randomBigInt, err := rand.Int(rand.Reader, new(big.Int).SetUint64(^uint64(0)))
	if err != nil {
		return 0
	}

	randomUint := randomBigInt.Uint64()

	return uint(randomUint)
}

func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		ID:   randomUint(),
		Conn: conn,
		Send: make(chan []byte),
	}

	hub.Register <- client

	go client.Read()
	go client.Write()
}
