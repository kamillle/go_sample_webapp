// 1つのクライアントへの接続を管理するclient typeの定義

package main

import (
	"github.com/gorilla/websocket"
)

type struct {
	// clientのためのwebsocket
	socket *websocket.Conn
	// メッセージが送られるチャネル
	send chan []byte
	// clientが参加しているチャットルーム
	room *room
}
