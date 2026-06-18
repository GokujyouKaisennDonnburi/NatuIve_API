// Package middleware は Gin のミドルウェアを提供する。
//
// 認証(Supabase が発行する JWT を JWKS で検証)・リクエストロギング・CORS などを
// ここに置く。最初の認証付きエンドポイントを追加するタイミングで実装する。
package middleware
