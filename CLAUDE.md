# Ghoi（語彙）

日本語・英語・韓国語の語彙を**引いて、溜めて、出会い直す**ためのアプリ。まずは自分ひとりで使う。

引くだけでは語彙力は上がらない。上がるのは同じ語に再会したときなので、`引く → 溜まる → 出会い直す` を中心の輪に置く。

## コマンド

道具とタスクは [mise](https://mise.jdx.dev) に集めている（`mise.toml`）。

| 目的 | コマンド |
|---|---|
| 道具を揃える | `mise install` |
| 使えるタスク一覧 | `mise tasks` |
| CI と同じ検査 | `mise run check` |
| テスト（競合検出つき） | `mise run test` |
| ビルド | `mise run build` → `bin/ghoi` |
| 整形 | `mise run fmt` |
| Terraform の検査 | `mise run tf-check`（GCP 接続不要） |
| Terraform の差分 | `mise run tf-plan`（ADC が要る） |

`mise run <task>` は `mise r <task>` と短く書ける。シェルに
`eval "$(mise activate zsh)"` を入れておくと、このディレクトリに入るだけで
`mise.toml` の Go に切り替わる（入れなくても `mise run` は正しい版を使う）。

## 構成

```
cmd/ghoi/            サーバ
internal/            アプリ本体（api / generate / srs / store / config）
db/migrations/       スキーマ
web/                 Vite + TanStack。package.json はここだけ
deploy/terraform/    GCP の構成
```

現在あるのは `cmd/ghoi` だけ。**残りは必要になった PR で作る。空のディレクトリを先に置かない。**

`mise.toml` の `[tools]` も同じ方針で、必要になった PR で足す（`web/` を作る PR で
node と pnpm、インフラの PR で terraform）。

## GCP

| | |
|---|---|
| プロジェクト ID | `ghoi-507101`（名前は `ghoi`） |
| リージョン | `asia-northeast1` |
| tfstate | `gs://ghoi-507101-tfstate`（バージョニング有効） |
| イメージ | `asia-northeast1-docker.pkg.dev/ghoi-507101/ghoi/ghoi` |
| Cloud Run | サービス名 `ghoi`、サービスアカウント `ghoi-run@...` |

**手元から触るには ADC が要る。**

```sh
gcloud auth application-default login
gcloud auth application-default set-quota-project ghoi-507101
```

`gcloud auth login` とは別物で、保存先も用途も違う（前者は gcloud コマンド自身、
後者は Terraform などのライブラリ）。確かめるには:

```sh
gcloud auth application-default print-access-token   # 出力は認証情報なので人に見せない
```

### state バケットだけは Terraform の外で作る

state を置く場所が無いと `terraform init` ができないため、ここだけ鶏と卵になる。
作り直すときは一度だけ次を実行する。

```sh
gcloud storage buckets create gs://ghoi-507101-tfstate \
  --project=ghoi-507101 --location=asia-northeast1 \
  --uniform-bucket-level-access --public-access-prevention
gcloud storage buckets update gs://ghoi-507101-tfstate --versioning
```

**バージョニングは必須。** state を壊しても前の版に戻せる。

### デプロイの順番

Cloud Run は**存在するイメージしか受け付けない**ので、初回だけ順番の制約がある。

```sh
mise run gcp-enable-apis   # 初回だけ。API を先に有効化する
mise run gcp-build         # Cloud Build がイメージを作って push する
mise run gcp-deploy        # そのイメージで Cloud Run を更新する
```

2回目以降は `gcp-build` → `gcp-deploy` の2つだけ。

**イメージのタグは git の短い SHA。** `latest` を使わないのは、
どのコミットが動いているのかが分からなくなるため。
ロールバックも「前の SHA で deploy し直す」でできる。

**手元に Docker は要らない。** `gcloud builds submit` がソースを送り、
Cloud Build がクラウド側で Dockerfile を使ってビルドする。

### `apply` は手動

CI は `plan` も `apply` もしない。`mise run tf-check`（整形と構文）だけを見る。
インフラが push で勝手に変わらないほうが安心だし、`plan` を読む練習にもなる。

```sh
cd deploy/terraform && terraform apply
```

## 進め方

**1 PR = 1テーマ。**「ついでに」を入れない。マージ後もリポジトリは常に動く状態を保つ。

1. 対話で、次の PR に何を含めるかを決める（設計の判断もここで済ませる）
2. ブランチを切って実装する。テストも同じ PR に入れる
3. PR の本文に「何を・なぜそうしたか」を書く
4. CI が通ったらマージし、**学んだことをこのファイルに足す**

差分は 200〜400 行が目安。読み返して意味が分かる規模を超えたら分ける。

全体の計画（21 PR / 8フェーズ）は別途まとめてある。

## 決まっていること

### 語ではなく語義（sense）を単位にする

`hot` を1件として持つと「暑い / 辛い / 人気の」が混ざり、復習として機能しない。
**保存も復習のスケジュールも sense 単位。** ここを外すと後で全部やり直しになる。

### 生成結果は恒久キャッシュする

例文と訳は LLM で生成し、Postgres に永久に保存する。同じ語は二度と生成しない。
**費用は「使う回数」ではなく「語彙の数」で決まる。**

`generation` テーブルに**モデル名とプロンプトの版**を記録する。これが無いと、
あとでプロンプトを良くしたときに「どれを作り直すべきか」が分からなくなる。

### `review_log` は最初から全部残す

復習アルゴリズムは SM-2 で始めるが、**ログさえ残っていれば後から FSRS に載せ替えられる。**
捨てると乗り換えられない。集計値だけでなく1回ごとの記録を持つ。

### 生成元は差し替えられるようにする

`Generator` インターフェースを切り、実装を Gemini / Claude / OpenAI で並べる。
まず Gemini の無料枠で作り、**実際に同じ語を複数モデルで生成して読み比べてから**本番のモデルを決める。
韓国語の語感が合っているかは、読んで判断するしかない。

### フロントは Go に同梱する

`go:embed` で SPA をバイナリに埋め込む。成果物1つ、デプロイ1回、CORS なし。
フロントだけ頻繁に出したくなったら、そのとき分ける。

### Terraform は同居させるが、自動 `apply` はしない

PR では `plan` を出すだけ。`apply` は手動。インフラが勝手に変わらないほうが安心だし、
`plan` を読む練習にもなる。

`.terraform.lock.hcl` は**コミットする**。しないと人や CI ごとに違う provider の版が入る。

### Cloud Run のサービスアカウントは最小から

既定のサービスアカウントは権限が広すぎるので、`ghoi-run` を専用に作ってある。
**いまは権限をひとつも付けていない。** Gemini を叩く PR、Cloud SQL に繋ぐ PR で、
必要な role をその都度足す。

### 公開範囲

Cloud Run は既定で非公開。`allUsers` に `roles/run.invoker` を付けて公開している。
**いま公開されているのは `/healthz` だけ。** 語を引く API を足す前に認証（PR 14-15）を入れること。

### API は使う PR で有効にする

GCP は既定でほとんどの API が無効なので、使う前に明示的に有効化する。
`apis.tf` に並べるのは**その時点で実際に使うもの**だけにして、先回りして足さない。

`disable_on_destroy = false` にしてある。無効化すると、その API で作った資源が
壊れることがあるため。

### JS の workspace ツールは入れない

pnpm workspace や Turborepo が効くのは、互いに依存する JS パッケージが複数あるとき。
**Ghoi の JS は `web/` 1つだけ**で、ビルドの本体は Go 側にある。

### タスクは mise に集める。Makefile は使わない

Go と Node を跨ぐので、言語専用のタスクランナー（`package.json` の scripts など）は上位に置けない。
mise を選んだのは、**タスクだけでなく言語のバージョンも同じファイルで固定できる**から。
CI でも `jdx/mise-action` で同じ `[tools]` を入れるので、**手元と CI の版が必ず一致する。**

Makefile も候補だったが、タブ必須・`.PHONY`・`$$` のエスケープ・ヘルプの自作といった
「道具と戦う」部分が多く、macOS 同梱の make は 3.81（2006年）で古い。

### コンテナは手元でビルドしない

`gcloud run deploy --source .` は **Dockerfile があればそれを使い、Cloud Build がクラウド側でビルドする**。
手元に Docker のエンジン（Colima / Docker Desktop など）は要らない。macOS で Docker を動かすには
Linux VM が要るので、まだ必要と分かっていないうちは常駐させない。

Dockerfile の間違いは Cloud Build のログに出る。デバッグが実際に苦痛になったら、そのとき Colima を入れる。

### Cloud Run の約束事

- **`PORT` 環境変数を読む。** 既定は 8080 だが、渡された値に従う
- **ホストを絞らずに待つ**（`:8080`）。`127.0.0.1:8080` にするとコンテナの外から届かず、起動に失敗する
- **SIGTERM を受けたら畳む。** Cloud Run は終了時に SIGTERM を送り、しばらく待ってから止める。
  無視すると処理中のリクエストが切れる
- `/healthz` で DB や外部 API を確かめない。**アプリが待ち受けを始めたかどうか**だけを返す

## 踏みやすいところ

- **`//go:embed` のパスは `..` を使えない。** そのソースファイルのあるディレクトリ基準なので、
  `internal/api/` から `web/dist` は埋め込めない。`web/embed.go`（`package web`）を置く
- **`//go:embed` はパターンが1つもマッチしないとビルドエラーになる。**
  クローン直後は `web/dist` が無いので `web/dist/.gitkeep` を置いておく
- **`go run` はシグナルを子プロセスに渡さない。** `mise run dev` に SIGTERM を送っても、
  終了処理は走らず 143 で殺される。**畳む動きを確かめるときはビルドしたバイナリを直接動かす**
- **`CGO_ENABLED=0` を忘れると distroless で動かない。** libc に依存したバイナリになるが、
  distroless には libc が無い
- **`scratch` ではなく `distroless/static` を使う。** CA 証明書が入っているため。
  無いと Gemini などへの HTTPS が証明書検証で落ちる
- **mise を入れただけでは PATH は切り替わらない。** `mise activate` をシェルに入れていない場合、
  素の `go` は Homebrew のものが使われる。`mise.toml` の版で動かしたいときは `mise run` か `mise exec` を通す

## 秘密の扱い

- **API キーをリポジトリに置かない。** 手元は `.env`（`.gitignore` 済み）、本番は Secret Manager
- GCP へのデプロイは **Workload Identity Federation** を使い、長期のサービスアカウントキーを作らない
- Gemini の無料枠は入力が製品改善に使われる規約になっている。**他人のデータを入れない**

## 費用

- Cloud Run はゼロスケールなので、リクエストが無ければほぼ無料
- **月額が立ち上がるのは Cloud SQL を作った時点だけ。** 停止と再開の手順を README に書くこと
- 生成の費用は語彙の数で決まる。1万語ためても、モデルによるが総額 $30〜$220 の範囲
