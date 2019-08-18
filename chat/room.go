// クライアントとの接続管理とメッセージのルーティングを持つ

package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// clientsの在室を管理する clients map は直接操作しない
// clientの入退出は必ず join と leave チャネルを使って管理する
// clientsのmapを直接操作するとメモリの破壊等が起きる可能性がある
type room struct {
	// 他のclientに転送するためのメッセージを保持するチャネル
	// メッセージを受け取ったら全てのclientに対してメッセージを送信する
	forward chan []byte
	// ルーム(clients)に参加するためのチャネル
	join chan *client
	// ルーム(clients)から退室するためのチャネル
	leave chan *client
	// room内にいる全clientsが保持される
	clients map[*client]bool
}

// チャットルームを生成するメソッド
// NOTE: room型への定義ではないので注意
func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) run() {
	// ループ処理。強制終了しない限り回り続ける。
	// goroutine(バックグラウンド)で実行するため、他の処理がブロックされることはない
	//
	// ループ処理の中で join, leave, forward のチャネルを監視しており、メッセージが届くとcase文が評価される
	for {
		select {
		case client := <-r.join:
			// 参加
			r.clients[client] = true
		case client := <-r.leave:
			// 退出
			// r.clients から client を削除(退出)する
			delete(r.clients, client)
			// client.send チャネルを閉じることで以後のメッセージの受信をしないことになる
			close(client.send)
		case msg := <-r.forward:
			// 全てのclient(r.clients)にメッセージを送信する
			for client := range r.clients {
				select {
				case client.send <- msg:
					// メッセージ送信
				default:
					// メッセージ送信に失敗
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

// HTTPハンドラとしての機能を実装する
//
// websocketを利用するには、websocket.Upgrader型を利用し、HTTP通信をupgradeする必要がある
//   ref: https://qiita.com/south37/items/6f92d4268fe676347160#1-%E3%82%B3%E3%83%8D%E3%82%AF%E3%82%B7%E3%83%A7%E3%83%B3%E7%A2%BA%E7%AB%8B
//
func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)

	// websocketコネクションの確立に失敗したらreturnする
	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	// 生成したクライアントをroomに入れる
	r.join <- client

	// client.readがループなのでこのdeferはクライアントの退出時まで評価されない
	// クライアントが終了を支持するまでこのdeferは遅延される
	defer func() { r.leave <- client }()
	// goroutineを立ち上げてスレッドによる並列処理を行う
	go client.write()
	client.read()
}
