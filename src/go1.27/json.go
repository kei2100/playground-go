package go1_27

// encoding/json　v1, v2 のデフォルト挙動の代表的な違い
//
// - 不正なUTF-8
//   - v1: U+FFFD(�) に置換
//   - v2: エラー
//   - 戻すオプション: jsontext.AllowInvalidUTF8(true)
//- 重複キー
//   - v1: 許容（後の値が前の値を置換またはマージ）
//   - v2: エラー
//   - 戻すオプション: jsontext.AllowDuplicateNames(true)
//- フィールド名マッチ
//   - v1: 大文字小文字を無視
//   - v2: 完全一致のみ
//   - 戻すオプション: MatchCaseInsensitiveNames(true)
//- nilスライス/マップ
//   - v1: null
//   - v2: 通常は [] / {}
//   - 戻すオプション: FormatNilSliceAsNull(true) / FormatNilMapAsNull(true)
//- time.Duration
//   - v1: ナノ秒の数値
//   - v2: エラー
//   - 戻すオプション: FormatDurationAsNano(true)
