package model

// HealthResponse はヘルスチェックのレスポンス。
type HealthResponse struct {
	// Status はサーバの稼働状態。正常時は "ok"。
	Status string `json:"status" example:"ok"`
}
