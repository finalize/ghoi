package config

import (
	"strings"
	"testing"
)

// env はテスト用の getenv を作る。
func env(pairs map[string]string) func(string) string {
	return func(k string) string { return pairs[k] }
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		environ  map[string]string
		wantAddr string
	}{
		{"PORT が無ければ 8080", nil, ":8080"},
		{"PORT が空でも 8080", map[string]string{"PORT": ""}, ":8080"},
		{"Cloud Run が渡してくる値を使う", map[string]string{"PORT": "8080"}, ":8080"},
		{"別の番号も通る", map[string]string{"PORT": "3000"}, ":3000"},
		{"上限も通る", map[string]string{"PORT": "65535"}, ":65535"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Load(env(tt.environ))
			if err != nil {
				t.Fatalf("エラー: %v", err)
			}
			if got.Addr != tt.wantAddr {
				t.Errorf("Addr = %q, 期待は %q", got.Addr, tt.wantAddr)
			}
		})
	}
}

// ホスト名を付けないこと。127.0.0.1 に絞るとコンテナの外から届かない。
func TestLoadListensOnAllInterfaces(t *testing.T) {
	got, err := Load(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(got.Addr, ":") {
		t.Errorf("Addr = %q, ホストを絞らずに待つこと", got.Addr)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name string
		port string
		want string
	}{
		{"数値ではない", "abc", "数値ではありません"},
		{"空白入り", "80 80", "数値ではありません"},
		{"0 は無効", "0", "範囲外"},
		{"負の数", "-1", "範囲外"},
		{"65536 は無効", "65536", "範囲外"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(env(map[string]string{"PORT": tt.port}))
			if err == nil {
				t.Fatalf("PORT=%q はエラーになるはずです", tt.port)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Errorf("エラー = %v, %q を含むはず", err, tt.want)
			}
		})
	}
}
