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
	// ケース1: 素直な足し算（まだ mod は跨がない）
	got := feAdd(big.NewInt(2), big.NewInt(3))
	if got.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("feAdd(2,3) = %v, want 5", got)
	}

	// ケース2: p を跨ぐ足し算。 (p-1) + 3 = p + 2 ≡ 2 (mod p)
	pm1 := new(big.Int).Sub(p, big.NewInt(1)) // p-1
	got = feAdd(pm1, big.NewInt(3))
	if got.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("feAdd(p-1,3) = %v, want 2", got)
	}
}

// TODO: TestFeSub を自分で書く。ヒント: 3 - 5 は 0 未満。mod p で正の値に畳めているか？
// TODO: TestFeMul を自分で書く。ヒント: (p-1) * 2 ≡ p-2 (mod p) など。
// TODO: TestFeInv を自分で書く。property test: feMul(a, feInv(a)) == 1。

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
