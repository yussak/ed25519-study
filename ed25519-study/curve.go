package main

// curve.go : ねじれ Edwards 曲線
//   -x^2 + y^2 = 1 + d*x^2*y^2   (mod p)
// 上の点と、その群演算(加算・スカラー倍)を実装する場所。
//
// field.go の feAdd/feSub/feMul/feInv が動いてから取り組むと良い。
//
// --- 用意したいもの ---
//
// TODO(1): 曲線定数 d = -121665 / 121666 (mod p) を定義する。
//   ヒント: feMul(-121665, feInv(121666)) の要領。負の数は mod p で正に。
//
// TODO(2): 点の表現を決める。まずは分かりやすさ優先でアフィン座標 (x, y) から。
//   type Point struct { X, Y *big.Int }
//   (高速化したくなったら拡張座標へ。最初はやらない。)
//
// TODO(3): 基点(base point) B を定義する。RFC 8032 の既知の値を使う。
//
// TODO(4): 点の加算 add(P, Q) を Edwards 加算公式で実装する。
//
// TODO(5): スカラー倍 scalarMul(k, P) を実装する(まずは素直な double-and-add)。
//
// 【これが出れば成功 (curve 単体)】
//   ・base 点 B が曲線の式を満たす (onCurve(B) == true)。
//   ・add(B, B) と scalarMul(2, B) が一致する。
//   ・scalarMul(L, B) が単位元になる (L は RFC 8032 の群位数)。※応用課題
