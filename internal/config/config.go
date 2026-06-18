package config

import (
	"os"
	"strings"
)

// Config はアプリ全体の設定値を保持する。環境変数から読み込む。
type Config struct {
	// Port はサーバの待ち受けポート。
	Port string
	// TrustedProxies は信頼するプロキシ(CIDR/IP)。nil = どのプロキシも信頼しない。
	TrustedProxies []string
}

// Load は環境変数から Config を構築する。
//
// PORT が未設定なら 8080、TRUSTED_PROXIES が未設定なら nil(直接接続のみ信頼)。
func Load() Config {
	cfg := Config{
		Port: os.Getenv("PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if v := os.Getenv("TRUSTED_PROXIES"); v != "" {
		cfg.TrustedProxies = strings.Split(v, ",")
	}
	return cfg
}
