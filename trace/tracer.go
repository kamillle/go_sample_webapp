package trace

import (
	"fmt"
	"io"
)

// 大文字にしておくことで外部アクセスが可能になる。小文字はプライベート
type Tracer interface {
	Trace(...interface{})
}

// ユーザーはこのAPIを利用して、tracerオブジェクトを受け取ることが可能
// tracerオブジェクトはプライベートなため直接操作はできないが、パブリックなAPIを利用して
// オブジェクトを受け取ることは可能
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// こちらは小文字からスタートなのでプライベートな構造体となる
type tracer struct {
	out io.Writer
}

// ...interface{} 型は任意の方の引数を0~複数個受け取れる
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
