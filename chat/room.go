// クライアントとの接続管理とメッセージのルーティングを持つ

package main

type room struct {
	// 他のclientに転送するためのメッセージを保持するチャネル
	forward chan byte[]
}
