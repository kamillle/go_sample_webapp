package main

import (
	"flag"
	"log"
	"net/http"
	_ "os"
	"path/filepath"
	"sync"
	"text/template"

	_ "github.com/kamillle/go_sample_webapp/trace"
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
	// rを渡すことでテンプレート側でリクエスト情報を参照できるようにしている
	t.templ.Execute(w, r) // t.templ.Executeの戻り地はチェックすべきらしい
}

func main() {
	// flagパッケージはコマンドライン引数を扱えるようにする
	// flag.String(<パラメータ名>, <デフォルト値>, <パラメータの説明>)
	// flag.Stringは *string型(フラグの値が保持されているアドレス)を返すため、値自体を参照したい場合は `*` 関節演算子を利用する
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// flag.Parse() でコマンドラインのパラメータをパースし、変数への値の代入処理が行われる
	flag.Parse()

	room := newRoom()
	// コメントインするとロギングが行われる
	// room.tracer = trace.New(os.Stdout)

	// templateHnadler型のオブジェクトを生成して、そのアドレスを渡している
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	// HTTPハンドラが実装されたroom型のオブジェクトを渡す
	http.Handle("/room", room)

	// ループ処理の中でwebsocket通信を利用する
	go room.run()

	log.Println("Webサーバーを開始します。ポート: ", *addr)

	// start server
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
