package service

import (
	"fmt"
	"strings"
)

const (
	// maxFilenameLen はファイル名の最大長（rune 単位）。DB カラムと表示の暴走を防ぐ。
	maxFilenameLen = 255
	// fallbackFilename はファイル名が空・不明なときのダウンロード名。
	fallbackFilename = "download"
)

// sanitizeFilename はユーザー提供のファイル名を安全な保存・表示用の値へ正規化する。
//
// - パス区切り（/ や \）を除去し basename のみにする（パストラバーサル防止）。
// - 制御文字（CR/LF/タブ含む）を除去する（Content-Disposition ヘッダインジェクション防止）。
// - 前後空白を除去し、rune 単位で maxFilenameLen に切り詰める。
//
// 正規化の結果が空になった場合は空文字を返す（フォールバックは呼び出し側の責務）。
func sanitizeFilename(name string) string {
	if i := strings.LastIndexAny(name, `/\`); i >= 0 {
		name = name[i+1:]
	}
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, name)
	name = strings.TrimSpace(name)
	if runes := []rune(name); len(runes) > maxFilenameLen {
		name = strings.TrimSpace(string(runes[:maxFilenameLen]))
	}
	return name
}

// contentDispositionInline は inline 表示用の Content-Disposition ヘッダ値を生成する。
//
// inline とすることでブラウザ内プレビュー（PDF など）を維持しつつ、保存時の
// 既定ファイル名を元名にできる。日本語等の非 ASCII は RFC 5987 の filename*（UTF-8）
// で表現し、非対応クライアント向けに ASCII フォールバックの filename= も併記する。
func contentDispositionInline(filename string) string {
	name := sanitizeFilename(filename)
	if name == "" {
		name = fallbackFilename
	}
	return fmt.Sprintf("inline; filename=%q; filename*=UTF-8''%s", asciiFallback(name), rfc5987Encode(name))
}

// asciiFallback は非 ASCII・制御文字を '_' に置換した ASCII フォールバック名を返す。
// filename= 用（quoted-string としてクォートするのは呼び出し側の %q）。
func asciiFallback(name string) string {
	var b strings.Builder
	for _, r := range name {
		if r < 0x20 || r > 0x7e {
			b.WriteByte('_')
			continue
		}
		b.WriteRune(r)
	}
	s := b.String()
	if strings.TrimSpace(s) == "" {
		return fallbackFilename
	}
	return s
}

// rfc5987Encode は RFC 5987 の value-chars に従い UTF-8 バイト列をパーセントエンコードする。
func rfc5987Encode(s string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	for _, c := range []byte(s) {
		if isAttrChar(c) {
			b.WriteByte(c)
			continue
		}
		b.WriteByte('%')
		b.WriteByte(upperhex[c>>4])
		b.WriteByte(upperhex[c&0xf])
	}
	return b.String()
}

// isAttrChar は RFC 5987 の attr-char（パーセントエンコード不要な文字）かを判定する。
func isAttrChar(c byte) bool {
	switch {
	case c >= 'A' && c <= 'Z', c >= 'a' && c <= 'z', c >= '0' && c <= '9':
		return true
	}
	switch c {
	case '!', '#', '$', '&', '+', '-', '.', '^', '_', '`', '|', '~':
		return true
	}
	return false
}
