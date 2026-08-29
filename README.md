# Ghoi

日本語・英語・韓国語の語彙を**引いて、溜めて、出会い直す**ためのアプリ。

語を入れると、例文と他の2言語での言い方が出る。引いた語義は保存でき、間隔を空けて復習に出てくる。
引くだけでは語彙力は上がらない。上がるのは同じ語に再会したときなので、そこまでを1つの道具にする。

## 状態

**作りはじめたところ。** いまは骨格だけで、まだ何も提供しない。

| | |
|---|---|
| フロント | Vite + TanStack（未着手） |
| API | Go / Cloud Run |
| DB | PostgreSQL |
| 生成 | Gemini（無料枠から） |
| 構成管理 | Terraform |
| CI/CD | GitHub Actions + Workload Identity Federation |

## 開発

道具とタスクは [mise](https://mise.jdx.dev) で管理している。

```sh
brew install mise   # 初回だけ
mise install        # mise.toml に書かれた Go などを入れる

mise tasks          # 使えるタスクの一覧
mise run check      # CI と同じ検査（gofmt / vet / test）
mise run build      # bin/ghoi を作る
```

Go のバージョンは `mise.toml` で固定してあるので、手元と CI で同じものが使われる。

リポジトリの決まりごと・設計判断・踏みやすいところは [CLAUDE.md](CLAUDE.md) にまとめてある。

## ライセンス

MIT
