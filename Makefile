# Ghoi の開発用コマンド。Go と Node を跨ぐので、言語中立な Make にしている。
# 1つ増やすたびに CLAUDE.md の表も更新すること。

BIN     := bin/ghoi
PKG     := ./...
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -X main.version=$(VERSION)

.DEFAULT_GOAL := help
.PHONY: help fmt vet test build check clean

help: ## このヘルプを出す
	@grep -hE '^[a-z-]+:.*?## ' $(MAKEFILE_LIST) | awk -F':.*?## ' '{printf "  \033[36m%-8s\033[0m %s\n", $$1, $$2}'

fmt: ## gofmt をかける
	gofmt -w .

vet: ## go vet
	go vet $(PKG)

test: ## テスト（競合検出つき）
	go test -race $(PKG)

build: ## バイナリを作る
	go build -ldflags '$(LDFLAGS)' -o $(BIN) ./cmd/ghoi

check: vet test ## CI と同じ検査をまとめて走らせる
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then echo 'gofmt が必要:'; echo "$$unformatted"; exit 1; fi
	@echo 'check: ok'

clean: ## 生成物を消す
	rm -rf bin
