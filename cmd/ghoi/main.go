// ghoi は日英韓の語彙を引いて、溜めて、出会い直すためのアプリ。
//
// いまは骨格だけで、まだ何も提供しない。
// PR 2 で /healthz を返す HTTP サーバになる。
package main

import (
	"flag"
	"fmt"
	"os"
)

// version はビルド時に Makefile が埋め込む（-ldflags "-X main.version=...")。
// 手で go run したときは dev のまま。
var version = "dev"

func main() {
	showVersion := flag.Bool("version", false, "バージョンを出して終わる")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}

	fmt.Fprintf(os.Stderr, "ghoi %s — まだ骨格だけです\n", version)
}
