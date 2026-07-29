package main

// sign.go : Ed25519 の鍵生成・署名・検証 (RFC 8032 §5.1) を実装する場所。
//
// field.go と curve.go が動いてから最後に取り組む。
// ハッシュは crypto/sha512 を使う想定。
//
// --- 実装する関数の想定シグネチャ ---
//
// TODO(1): 鍵生成
//   func generateKey(seed []byte) (pub []byte, ...)
//   ・SHA-512(seed) の前半32byteをクランプしてスカラー a を作る
//   ・A = a * B をエンコードして公開鍵にする
//
// TODO(2): 点のエンコード/デコード (32byte, y座標 + x の符号ビット)
//   func encodePoint(P) []byte / func decodePoint([]byte) (Point, error)
//   ※ ここは RFC 8032 の仕様を丁寧に。詰まりやすいポイント。
//
// TODO(3): 署名 sign(priv, msg) []byte  (R || S の 64byte)
//   ・r = SHA-512(prefix || msg) を法 L で還元
//   ・R = r * B
//   ・k = SHA-512(encode(R) || encode(A) || msg) を法 L で還元
//   ・S = (r + k*a) mod L
//
// TODO(4): 検証 verify(pub, msg, sig) bool
//   ・8*S*B == 8*R + 8*k*A  を確認する
//
// 【これが出れば成功 (最終ゴール)】
//   ・自分で作った鍵で sign → verify が true。
//   ・1bit でも msg を変えると verify が false。
//   ・RFC 8032 のテストベクタ(既知の seed/msg/署名)と完全一致する。
//     → crypto/ed25519 標準ライブラリの出力と突き合わせるのも良い検証法。
