package trace

// 大文字にしておくことで外部アクセスが可能になる。小文字はプライベート
type Tracer interface {
	// ...interface{} 型は任意の方の引数を0~複数個受け取れる
	Trace(...interface{})
}
