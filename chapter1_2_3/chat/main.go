package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	_ "os"
	"path/filepath"
	"sync"
	"text/template"

	_ "github.com/kamillle/go_sample_webapp/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	fmt.Println(data)
	// 独自に定義したdataを渡すことでテンプレート側で参照できるようになる
	t.templ.Execute(w, data) // t.templ.Executeの戻り地はチェックすべきらしい
}

func main() {
	// flagパッケージはコマンドライン引数を扱えるようにする
	// flag.String(<パラメータ名>, <デフォルト値>, <パラメータの説明>)
	// flag.Stringは *string型(フラグの値が保持されているアドレス)を返すため、値自体を参照したい場合は `*` 関節演算子を利用する
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	// flag.Parse() でコマンドラインのパラメータをパースし、変数への値の代入処理が行われる
	flag.Parse()
	// Gomniauthのセットアップ
	// SetSecurityKey の引数はクライアント・サーバー間処理の進行状態管理の際に行われるデジタル署名用
	// ここではランダムな値、または自分で決めた文字列でいい
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		facebook.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/facebook"),
		github.New("クライアントID", "秘密の値", "http://localhost:8080/auth/callback/github"),
		google.New(os.Getenv("GO_SAMPLE_WEBAPP_GOOGLE_CLIENT_ID"), os.Getenv("GO_SAMPLE_WEBAPP_GOOGLE_CLIENT_SECRET"), "http://localhost:8080/auth/callback/google"),
	)

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
