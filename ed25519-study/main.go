package main

import "fmt"

// ed25519-study : Ed25519 (RFC 8032) をゼロから自作して理解する。
// まずは 1ファイル(main.go)にまとめて進める。分割は中身が分かってきてから。
// 使うのは標準ライブラリのみ (math/big, crypto/sha512)。

func main() {
	// 実装が進んだら、ここで鍵生成 → 署名 → 検証のデモを動かす予定。
	fmt.Println("ed25519-study: hello")
}

// ============================================================
// 1. 有限体 GF(p), p = 2^255 - 19   （まずここから）
// ============================================================
// TODO(1): 素数 p = 2^255 - 19 を定義
//   ヒント: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(19))
// TODO(2): feAdd(a, b) = (a + b) mod p
// TODO(3): feSub(a, b) = (a - b) mod p   (負にならない mod に注意)
// TODO(4): feMul(a, b) = (a * b) mod p
// TODO(5): feInv(a) = a^(p-2) mod p       (フェルマーの小定理)
//   → 成功: feMul(a, feInv(a)) == 1

// ============================================================
// 2. ねじれ Edwards 曲線  -x^2 + y^2 = 1 + d*x^2*y^2 (mod p)
// ============================================================
// TODO: 曲線定数 d = -121665 / 121666 (mod p)
// TODO: 点の表現 (まずアフィン (x, y))、基点 B、onCurve(P)
// TODO: add(P, Q) / scalarMul(k, P)
//   → 成功: onCurve(B) かつ add(B,B) == scalarMul(2,B)

// ============================================================
// 3. 鍵生成 / 署名 / 検証 (RFC 8032 §5.1)
// ============================================================
// TODO: encodePoint / decodePoint (32byte)
// TODO: generateKey(seed) / sign(priv, msg) / verify(pub, msg, sig)
//   → 成功: 自作鍵で sign→verify が true、1bit変えると false、
//           最後に RFC 8032 テストベクタ / 標準 crypto/ed25519 と一致
