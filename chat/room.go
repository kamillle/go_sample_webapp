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
