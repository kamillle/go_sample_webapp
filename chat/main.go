package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type templateHandler struct {
	once     sync.Once // sync.Once型を使用してテンプレートを一度だけコンパイルする
	filename string
	templ    *template.Template // 1つのテンプレートを表します
}

// templateHandler型へのメソッド定義
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// ./chat/. でgo runを実行すること
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r) // t.templ.Executeの戻り地はチェックすべきらしい
}

func main() {
	room := newRoom()

	// templateHnadler型のオブジェクトを生成して、そのアドレスを渡している
	http.Handle("/", &templateHandler{filename: "chat.html"})
	// HTTPハンドラが実装されたroom型のオブジェクトを渡す
	http.Handle("/room", room)

	// ループ処理の中でwebsocket通信を利用する
	go room.run()

	// start server
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
