// 1つのクライアントへの接続を管理するclient typeの定義

package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type client struct {
	// clientのためのwebsocket
	socket *websocket.Conn
	// メッセージが送られるチャネル
	// 送信者や送信時間をchatのやり取りに加えるために独自定義のmessage型を使っている
	send chan *message
	// clientが参加しているチャットルーム
	room *room
	// ユーザーに関する情報を保持する
	userData map[string]interface{}
}

func (c *client) read() {
	for {
		var msg *message
		// c.socket.ReadJSON()で message型オブジェクトのデコードの結果、errにオブジェクトが入らなければtrue
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
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
		// c.socket.WriteMessage の結果、 errがnilであればtrue
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
