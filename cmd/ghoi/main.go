// ghoi は日英韓の語彙を引いて、溜めて、出会い直すためのアプリ。
//
// いまは /healthz を返すだけ。語を引く API は PR 9 で足す。
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/finalize/ghoi/internal/api"
	"github.com/finalize/ghoi/internal/config"
)

// version はビルド時に埋め込む（-ldflags "-X main.version=..."）。
// 手で go run したときは dev のまま。
var version = "dev"

func main() {
	// 実処理を run に分けているのは、main が error を返せないため。
	// getenv と stderr を引数にしておくと、あとでテストから呼べる。
	ctx := context.Background()
	if err := run(ctx, os.Getenv, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "ghoi:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, getenv func(string) string, stderr io.Writer) error {
	showVersion := flag.Bool("version", false, "バージョンを出して終わる")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return nil
	}

	cfg, err := config.Load(getenv)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: api.New(),
		// ヘッダを送りきらない接続に居座られないようにする。
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Cloud Run は終了時に SIGTERM を送り、しばらく待ってから止める。
	// 受け取ったら処理中のリクエストを返しきってから閉じる。
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ListenAndServe は塞がるので、別の流れで走らせる。
	// 「起動に失敗した」と「終了しろと言われた」のどちらが先に来ても
	// 拾えるように、channel と select でまとめる。
	errCh := make(chan error, 1)
	go func() {
		fmt.Fprintf(stderr, "ghoi %s: %s で待ち受けます\n", version, cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		fmt.Fprintln(stderr, "ghoi: 終了します")
	}

	// 閉じるほうにも制限時間を付ける。ここで ctx を使い回すと
	// もう終わっている ctx を渡すことになり、待たずに切れてしまう。
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutCtx)
}
