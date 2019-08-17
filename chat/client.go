// 1つのクライアントへの接続を管理するclient typeの定義

package main

import (
	"github.com/gorilla/websocket"
)

type client struct {
	// clientのためのwebsocket
	socket *websocket.Conn
	// メッセージが送られるチャネル
	send chan []byte
	// clientが参加しているチャットルーム
	room *room
}

func (c *client) read() {
	for {
		if _, msg, err := c.socket.ReadMessage(); err == nil {
			// <- は左辺のチャネルに右辺の値を送信する演算子
			// = <- は右辺のチャネルから値を受信できる
			c.room.forward <- msg
		} else {
			// ループしたときなどに抜けられるようしておく
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
