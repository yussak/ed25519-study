package main

// field.go : 有限体 GF(p), p = 2^255 - 19 上の演算を実装する場所。
//
// Ed25519 の座標や中間計算はすべて mod p の世界で行う。
// まずはここを固めると、後段(curve/sign)が一気に楽になる。
//
// 実装の想定: math/big の *big.Int を使う。
//
// --- 最初に用意したいもの ---
//
// TODO(1): 素数 p = 2^255 - 19 を表す変数を定義する。
//   ヒント: new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 255), big.NewInt(19))
//
// TODO(2): 体の加算 feAdd(a, b) = (a + b) mod p
// TODO(3): 体の減算 feSub(a, b) = (a - b) mod p   (負にならないよう mod に注意)
// TODO(4): 体の乗算 feMul(a, b) = (a * b) mod p
// TODO(5): 体の逆元 feInv(a) = a^(p-2) mod p       (フェルマーの小定理)
//   ヒント: new(big.Int).Exp(a, pMinus2, p)
//
// 【これが出れば成功 (field 単体)】
//   feMul(a, feInv(a)) == 1  が任意の a(≠0) で成り立つ。
//   → field_test.go のテストを緑にすることがこの段の合格ライン。
