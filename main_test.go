package main

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/sha512"
	"math/big"
	"testing"
)

// テストは実装(main.go)と同じく 1ファイルにまとめる。分割は後で。

// ============================================================
// 内側ループ: 有限体などの小さな積み上げテスト
// ============================================================

func TestFeAdd(t *testing.T) {
	// ケース1: 和が p 未満なので mod しても変わらない
	got := feAdd(big.NewInt(2), big.NewInt(3))
	if got.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("feAdd(2,3) = %v, want 5", got)
	}

	// ケース2: 和が p 以上になるので mod で折り返す。(p-1)+3 = p+2 ≡ 2 (mod p)
	pm1 := new(big.Int).Sub(p, big.NewInt(1)) // p-1
	got = feAdd(pm1, big.NewInt(3))
	if got.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("feAdd(p-1,3) = %v, want 2", got)
	}
}

func TestFeSub(t *testing.T) {
	// ケース1: 折り返さない引き算
	got := feSub(big.NewInt(5), big.NewInt(3))
	if got.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("feSub(5,3) = %v, want 2", got)
	}

	// ケース2: 負になる引き算。3 - 5 = -2 ≡ p-2 (mod p)
	want := new(big.Int).Sub(p, big.NewInt(2)) // p-2
	got = feSub(big.NewInt(3), big.NewInt(5))
	if got.Cmp(want) != 0 {
		t.Errorf("feSub(3,5) = %v, want p-2", got)
	}
}

func TestFeMul(t *testing.T) {
	// ケース1: 折り返さない掛け算
	got := feMul(big.NewInt(2), big.NewInt(3))
	if got.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("feMul(2,3) = %v, want 6", got)
	}

	// ケース2: 折り返す掛け算。(p-1) * 2 = 2p-2 ≡ p-2 (mod p)
	pm1 := new(big.Int).Sub(p, big.NewInt(1))  // p-1
	want := new(big.Int).Sub(p, big.NewInt(2)) // p-2
	got = feMul(pm1, big.NewInt(2))
	if got.Cmp(want) != 0 {
		t.Errorf("feMul(p-1,2) = %v, want p-2", got)
	}
}

func TestFeInv(t *testing.T) {
	// feInv は mod p での逆元＝a に掛けて 1 になる整数（0.5 のような小数ではない）。
	// 返り値そのものを、手で導ける具体値と突き合わせて「特定の整数が返る」ことを示す。

	// ケース1: 2 の逆元は (p+1)/2。 2×(p+1)/2 = p+1 ≡ 1 (mod p) だから。
	half := new(big.Int).Div(new(big.Int).Add(p, big.NewInt(1)), big.NewInt(2)) // (p+1)/2
	if got := feInv(big.NewInt(2)); got.Cmp(half) != 0 {
		t.Errorf("feInv(2) = %v, want (p+1)/2 = %v", got, half)
	}

	// ケース2: p-1(≡ -1) の逆元は自分自身。 (-1)×(-1)=1 だから。
	pm1 := new(big.Int).Sub(p, big.NewInt(1))
	if got := feInv(pm1); got.Cmp(pm1) != 0 {
		t.Errorf("feInv(p-1) = %v, want p-1 = %v", got, pm1)
	}

	// ケース3: 1 の逆元は自分自身。
	if got := feInv(big.NewInt(1)); got.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("feInv(1) = %v, want 1", got)
	}

	// 定義（a × 逆元 ≡ 1）も確認。feInv と feMul の噛み合わせを見る。
	a := big.NewInt(2)
	if got := feMul(a, feInv(a)); got.Cmp(big.NewInt(1)) != 0 {
		t.Errorf("feMul(2, feInv(2)) = %v, want 1", got)
	}
}

func TestOnCurve(t *testing.T) {
	// ケース1: 基点 B は曲線上にある。
	// これが true になる = 定数 d と基点 B を正しく定義できた、の検算も兼ねる。
	if !onCurve(B) {
		t.Error("onCurve(B) = false, want true")
	}

	// ケース2: 単位元 (0,1) も曲線上。式に入れると -0+1 = 1+0 で成立する。
	id := point{X: big.NewInt(0), Y: big.NewInt(1)}
	if !onCurve(id) {
		t.Error("onCurve((0,1)) = false, want true")
	}

	// ケース3: 明らかに外れた点 (1,1) は曲線上にない。
	// 「常に true を返す」バグを弾くため、false になるべき点も確認する。
	off := point{X: big.NewInt(1), Y: big.NewInt(1)}
	if onCurve(off) {
		t.Error("onCurve((1,1)) = true, want false")
	}
}

func TestAdd(t *testing.T) {
	// ケース1: 単位元 (0,1) を足しても点は変わらない（P + O = P）。
	// 既に持っている単位元をオラクルに使い、既知の第3の点を用意せず縛る。
	id := point{X: big.NewInt(0), Y: big.NewInt(1)}
	if got := add(B, id); got.X.Cmp(B.X) != 0 || got.Y.Cmp(B.Y) != 0 {
		t.Errorf("add(B, id) = (%v, %v), want B", got.X, got.Y)
	}

	// ケース2: B の二重化 add(B, B) の結果も曲線上に留まる。
	// onCurve をオラクルにして「加算式が正しい」ことの強い証拠にする。
	if !onCurve(add(B, B)) {
		t.Error("onCurve(add(B, B)) = false, want true")
	}
}

// ============================================================
// 外側ループ: ゴールテスト（完成の定義、すべて t.Skip）
// 参照している API シグネチャは変更前提。
//
// ★ 当面のゴール = TestAgainstStdEd25519 を通すこと。
//    これの t.Skip を外して緑にできたら Ed25519 コア完成。
//    ph/ctx/batch はその後の応用（ゴールではない）。
// ============================================================

func fixedSeed() []byte {
	seed := make([]byte, ed25519.SeedSize)
	for i := range seed {
		seed[i] = byte(i)
	}
	return seed
}

// ★★★ 当面のゴール ★★★
// コア: 標準 crypto/ed25519 をオラクルにして keygen/sign/verify を突き合わせる。
// これを緑にするのが今の最終目標。
func TestAgainstStdEd25519(t *testing.T) {
	t.Skip("当面のゴール: 実装が揃ったら Skip を外す")

	seed := fixedSeed()
	msg := []byte("hello ed25519")

	// オラクル
	stdPriv := ed25519.NewKeyFromSeed(seed)
	stdPub := stdPriv.Public().(ed25519.PublicKey)
	stdSig := ed25519.Sign(stdPriv, msg)

	// 自作
	pub := generateKey(seed) // 32byte の公開鍵を返す想定
	sig := sign(seed, msg)   // 64byte の署名を返す想定

	if !bytes.Equal(pub, stdPub) {
		t.Errorf("public key mismatch\n got=%x\nwant=%x", pub, stdPub)
	}
	if !bytes.Equal(sig, stdSig) {
		t.Errorf("signature mismatch\n got=%x\nwant=%x", sig, stdSig)
	}
	if !verify(pub, msg, sig) {
		t.Error("verify rejected a valid signature")
	}
	bad := append([]byte(nil), msg...)
	bad[0] ^= 1
	if verify(pub, bad, sig) {
		t.Error("verify accepted a tampered message")
	}
}

// 応用: Ed25519ph（先に SHA-512 でハッシュしてから署名）。
func TestEd25519ph(t *testing.T) {
	t.Skip("ゴール(応用): Ed25519ph")

	seed := fixedSeed()
	msg := []byte("hello ed25519ph")
	digest := sha512.Sum512(msg)

	stdPriv := ed25519.NewKeyFromSeed(seed)
	stdPub := stdPriv.Public().(ed25519.PublicKey)
	stdSig, err := stdPriv.Sign(nil, digest[:], &ed25519.Options{Hash: crypto.SHA512})
	if err != nil {
		t.Fatal(err)
	}

	sig := signPh(seed, msg)
	if !bytes.Equal(sig, stdSig) {
		t.Errorf("ph signature mismatch\n got=%x\nwant=%x", sig, stdSig)
	}
	if !verifyPh(stdPub, msg, sig) {
		t.Error("verifyPh rejected a valid signature")
	}
}

// 応用: Ed25519ctx（コンテキスト文字列付き）。
func TestEd25519ctx(t *testing.T) {
	t.Skip("ゴール(応用): Ed25519ctx")

	seed := fixedSeed()
	msg := []byte("hello ctx")
	ctx := "example-context"

	stdPriv := ed25519.NewKeyFromSeed(seed)
	stdPub := stdPriv.Public().(ed25519.PublicKey)
	stdSig, err := stdPriv.Sign(nil, msg, &ed25519.Options{Context: ctx})
	if err != nil {
		t.Fatal(err)
	}

	sig := signCtx(seed, msg, ctx)
	if !bytes.Equal(sig, stdSig) {
		t.Errorf("ctx signature mismatch\n got=%x\nwant=%x", sig, stdSig)
	}
	if !verifyCtx(stdPub, msg, sig, ctx) {
		t.Error("verifyCtx rejected a valid signature")
	}
}

// 応用: バッチ検証（標準にオラクルが無いので自己整合で確認）。
func TestBatchVerify(t *testing.T) {
	t.Skip("ゴール(応用): バッチ検証")

	seed := fixedSeed()
	pub := generateKey(seed)

	var pubs, msgs, sigs [][]byte
	for i := 0; i < 3; i++ {
		m := []byte{byte(i)}
		pubs = append(pubs, pub)
		msgs = append(msgs, m)
		sigs = append(sigs, sign(seed, m))
	}

	if !verifyBatch(pubs, msgs, sigs) {
		t.Error("batch verify rejected an all-valid set")
	}
	sigs[1][0] ^= 1
	if verifyBatch(pubs, msgs, sigs) {
		t.Error("batch verify accepted a set containing a bad signature")
	}
}
