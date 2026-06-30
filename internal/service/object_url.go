package service

import "strings"

// PublicURLResolver はオブジェクトキーを公開URLへ変換する。
//
// 公開バケット（R2 パブリックバケット/カスタムドメイン）を前提に、
// baseURL + "/" + key を決定論的に組み立てる。これにより同じキーは常に
// 同じ URL になり、本文インライン画像など個数不定の参照にもそのまま使える。
//
// 本文に保存するのは object_key（完全URLではない）を推奨する。配信先が
// 変わってもベースURLの差し替えだけで追従でき、移行コストを抑えられる。
type PublicURLResolver struct {
	// baseURL は末尾スラッシュなしに正規化済みの配信ベースURL。
	// 空文字なら URL を生成しない（機能無効）。
	baseURL string
}

// NewPublicURLResolver は baseURL を正規化して PublicURLResolver を返す。
func NewPublicURLResolver(baseURL string) PublicURLResolver {
	return PublicURLResolver{baseURL: strings.TrimRight(baseURL, "/")}
}

// URL は単一キーの公開URLを返す。baseURL 未設定・key 空なら空文字を返す。
func (r PublicURLResolver) URL(key string) string {
	if r.baseURL == "" || key == "" {
		return ""
	}
	return r.baseURL + "/" + strings.TrimLeft(key, "/")
}

// URLs は複数キーの公開URL一覧を返す。
//
// nil/空入力でも空スライス（非 nil）を返し、JSON で null ではなく [] になるようにする。
// baseURL 未設定時は空スライスを返す。
func (r PublicURLResolver) URLs(keys []string) []string {
	urls := make([]string, 0, len(keys))
	for _, k := range keys {
		if u := r.URL(k); u != "" {
			urls = append(urls, u)
		}
	}
	return urls
}
