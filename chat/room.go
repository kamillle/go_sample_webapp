// クライアントとの接続管理とメッセージのルーティングを持つ

package main

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
					delete(r.clients, clinet)
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
