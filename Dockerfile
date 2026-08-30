# syntax=docker/dockerfile:1

# ---- ビルド ----
# Go のバージョンは mise.toml と揃える。片方だけ上げないこと。
FROM golang:1.27 AS build

WORKDIR /src

# 依存がまだ無いので、まとめて写している。
# go.sum ができたら「先に go.mod と go.sum だけ写して go mod download」に分けて、
# 依存が変わらないかぎりダウンロードを飛ばせるようにする（PR 7 あたり）。
COPY . .

ARG VERSION=dev

# CGO_ENABLED=0 で静的リンクにする。これをしないと libc に依存したバイナリになり、
# distroless（libc を持たない）で動かない。
#   -trimpath ビルドしたマシンの絶対パスをバイナリに残さない
#   -s -w     デバッグ情報を落として小さくする
RUN CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags "-s -w -X main.version=${VERSION}" \
      -o /out/ghoi ./cmd/ghoi

# ---- 実行 ----
# distroless にはシェルもパッケージマネージャも入っていない。
# 侵入されても道具が無いので、攻撃の足場になりにくい。
#
# scratch ではなく static を選んだのは CA 証明書が入っているため。
# 無いと Gemini などへの HTTPS が証明書検証で失敗する（PR 8 で効いてくる）。
# :nonroot タグは root 以外のユーザーで動く。
FROM gcr.io/distroless/static:nonroot

COPY --from=build /out/ghoi /ghoi

# Cloud Run は EXPOSE を見ないが、何番で待つのかを人に伝えるために書いておく。
# 実際の番号は環境変数 PORT で決まる。
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/ghoi"]
