// クライアントとの接続管理とメッセージのルーティングを持つ

package main

type room struct {
	// 他のclientに転送するためのメッセージを保持するチャネル
	// メッセージを受け取ったら全てのclientに対してメッセージを送信する
	forward chan []byte
}
