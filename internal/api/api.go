// Package api は HTTP のルーティングとハンドラを持つ。
//
// いまは /healthz だけ。語を引く API は PR 9 で足す。
package api

import "net/http"

// New はルーティング済みのハンドラを返す。
//
// Go 1.22 以降の ServeMux はメソッドとパス変数を書けるので、
// ルータのライブラリは入れない。
func New() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	return mux
}

// healthz は生きているかどうかだけを返す。
//
// Cloud Run はコンテナが待ち受けを始めたかを見るので、
// ここで DB や外部 API を確かめない。それらが落ちていても
// アプリ自体は起動しているし、確かめに行くと起動が遅くなる。
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("ok\n"))
}
