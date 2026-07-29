# Ed25519 自作 学習 TODO

RFC 8032 の Ed25519 を Go でゼロから実装する。`math/big` と `crypto/sha512` のみ使用。
コードは自分で書く（穴埋め形式）。各段は「これが出れば成功」を満たしたらチェックする。

進める順番: **field → curve → sign**（下の層が固まってから次へ）。

> ★ 当面のゴール: **`TestAgainstStdEd25519`** を通すこと。
> （標準 crypto/ed25519 と keygen/sign/verify がバイト一致 = Ed25519 コア完成）
> ph/ctx/batch は応用でゴールではない。

---

## 0. 環境構築
- [x] `go version` 確認（go1.24.7）
- [x] `ed25519-study/` 作成・`go mod init`
- [x] `main.go` / `field.go` / `curve.go` / `sign.go` / `field_test.go` の雛形
- [x] `go run .` と `go test ./...` が通る最小状態
- [ ] RFC 8032 本文を一度ざっと読む（§5.1 が Ed25519） … まず全体像の把握が目的

## 1. field.go — 有限体 GF(p), p = 2^255 - 19
- [ ] 素数 `p = 2^255 - 19` を定義
- [ ] `feAdd(a, b) = (a + b) mod p`
- [ ] `feSub(a, b) = (a - b) mod p`（負にならない mod に注意）
- [ ] `feMul(a, b) = (a * b) mod p`
- [ ] `feInv(a) = a^(p-2) mod p`（フェルマーの小定理）
- [ ] `field_test.go` の `t.Skip` を消して本テストにする
- [ ] **成功判定:** 任意の `a(≠0)` で `feMul(a, feInv(a)) == 1`（`go test` 緑）

## 2. curve.go — ねじれ Edwards 曲線 -x^2+y^2 = 1+d·x^2·y^2
- [ ] 曲線定数 `d = -121665 / 121666 (mod p)` を定義
- [ ] 点の表現を決める（まずアフィン `(x, y)` から） … 拡張座標にするかは後で調査
- [ ] 基点 `B` を RFC 8032 の既知値で定義
- [ ] `onCurve(P)` … 点が曲線式を満たすか判定
- [ ] 点の加算 `add(P, Q)`（Edwards 加算公式）
- [ ] スカラー倍 `scalarMul(k, P)`（まず素直な double-and-add）
- [ ] `curve_test.go` を作る
- [ ] **成功判定:** `onCurve(B) == true` かつ `add(B,B) == scalarMul(2,B)`
- [ ] （応用）群位数 `L` で `scalarMul(L, B)` が単位元になる … L の値と単位元の扱いを調査

## 3. sign.go — 鍵生成 / 署名 / 検証（RFC 8032 §5.1）
- [ ] 点のエンコード `encodePoint(P) []byte`（32byte, y + x の符号ビット）
- [ ] 点のデコード `decodePoint([]byte) (Point, error)` … 平方根の計算方法を調査（詰まりやすい）
- [ ] 鍵生成 `generateKey(seed)`（SHA-512 → クランプ → `A = a·B`）
- [ ] 署名 `sign(priv, msg)`（`R || S` の 64byte）
- [ ] 検証 `verify(pub, msg, sig) bool`（`8·S·B == 8·R + 8·k·A`）
- [ ] `sign_test.go` を作る
- [ ] **成功判定:** 自作鍵で `sign → verify` が true / msg を 1bit 変えると false

## 4. 仕上げ・検証の強化
- [ ] RFC 8032 のテストベクタ（既知の seed/msg/署名）と完全一致を確認
- [ ] 標準ライブラリ `crypto/ed25519` の出力と突き合わせる
- [ ] `main.go` を鍵生成→署名→検証のデモに差し替える
- [ ] （調査）タイミング攻撃など「学習実装ゆえ本番非推奨」な点を README にメモ

---

### まだ決めていない/要調査（タスク化しておく）
- [ ] 座標系: アフィンのまま進めるか、途中で拡張座標へ移すか判断する
- [ ] `feSub` などで負値を正に畳む書き方を統一する（ヘルパを作るか）
- [ ] テストの置き場所: パッケージ内 `_test.go` のままで十分か確認
