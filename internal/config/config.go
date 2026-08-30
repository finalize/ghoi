// Package config は環境変数から設定を読む。
//
// 読み方を関数の引数（getenv）にしているのは、テストで os.Setenv を触らずに済ませるため。
// 環境変数はプロセス全体で共有されるので、テストが並列に走ると壊れる。
package config

import (
	"fmt"
	"net"
	"strconv"
)

// DefaultPort は PORT が渡されなかったときに使う番号。
// Cloud Run は既定で 8080 を渡してくるので、それに合わせてある。
const DefaultPort = 8080

// Config はアプリ全体の設定。増えたらここに足す。
type Config struct {
	// Addr は net/http にそのまま渡せる待ち受けアドレス。
	Addr string
}

// Load は環境変数から設定を組み立てる。
//
// getenv には os.Getenv を渡す。テストでは好きな関数を渡せる。
func Load(getenv func(string) string) (Config, error) {
	port := DefaultPort
	if raw := getenv("PORT"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			return Config{}, fmt.Errorf("PORT が数値ではありません: %q", raw)
		}
		if n < 1 || n > 65535 {
			return Config{}, fmt.Errorf("PORT が範囲外です: %d", n)
		}
		port = n
	}

	// ホスト名を空にして「すべてのアドレスで待つ」にする。
	// 127.0.0.1 に絞ると、コンテナの外から届かず Cloud Run のヘルスチェックが落ちる。
	return Config{Addr: net.JoinHostPort("", strconv.Itoa(port))}, nil
}
