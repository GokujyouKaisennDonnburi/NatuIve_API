# NatuIve API 開発用タスク
# 使い方: `make` または `make help` で一覧を表示

# swag は go.mod の swaggo/swag と同じバージョンに固定する（生成物のブレを防ぐ）。
SWAG_VERSION := v1.16.6
SWAG ?= $(shell go env GOPATH)/bin/swag
SWAG_ENTRY := cmd/api/main.go
SWAG_OUT := ./api

.DEFAULT_GOAL := help

.PHONY: help
help: ## このヘルプを表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## ローカルでサーバを起動 (go run)
	go run ./cmd/api

.PHONY: build
build: ## ビルド確認
	go build ./...

.PHONY: test
test: ## テスト実行
	go test ./...

.PHONY: tidy
tidy: ## 依存の整理
	go mod tidy

.PHONY: fmt
fmt: ## フォーマット
	go fmt ./...

.PHONY: vet
vet: ## go vet
	go vet ./...

.PHONY: swag-install
swag-install: ## swag CLI をインストール (バージョン固定)
	go install github.com/swaggo/swag/cmd/swag@$(SWAG_VERSION)

.PHONY: swag
swag: ## OpenAPI ドキュメントを生成 (docs/)
	$(SWAG) init -g $(SWAG_ENTRY) -o $(SWAG_OUT) --parseDependency --parseInternal

.PHONY: swag-check
swag-check: swag ## docs/ が最新か確認 (CI 用: 差分があれば失敗)
	@git diff --exit-code $(SWAG_OUT) || (echo "docs/ が古いです。'make swag' を実行してコミットしてください" && exit 1)

.PHONY: up
up: ## 開発用コンテナを起動 (docker compose)
	docker compose up

.PHONY: down
down: ## 開発用コンテナを停止
	docker compose down
