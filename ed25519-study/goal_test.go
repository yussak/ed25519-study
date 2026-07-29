package main

import (
	"bytes"
	"crypto/ed25519"
	"testing"
)

// goal_test.go : ゴールテスト（完成の定義）。
// アウトサイドイン TDD の外側ループ。実装が揃うまで t.Skip で寝かせておく。
//
// 標準ライブラリ crypto/ed25519 をオラクル(正解)にして、自作実装が
// バイト単位で一致することを確認する。ここが緑になったら完成。
//
// 参照している generateKey / sign / verify のシグネチャ = 目指す最終 API。
// （変更前提。実装しながら形が変わってよい）

func TestAgainstStdEd25519(t *testing.T) {
	t.Skip("ゴール: 実装が揃ったら Skip を外す")

	// 固定 seed（32byte）。中身は何でもよいが決め打ちにして再現性を持たせる。
	seed := make([]byte, ed25519.SeedSize)
	for i := range seed {
		seed[i] = byte(i)
	}
	msg := []byte("hello ed25519")

	// --- オラクル: 標準ライブラリ ---
	stdPriv := ed25519.NewKeyFromSeed(seed)
	stdPub := stdPriv.Public().(ed25519.PublicKey)
	stdSig := ed25519.Sign(stdPriv, msg)

	// --- 自作実装 ---
	pub := generateKey(seed) // 32byte の公開鍵を返す想定
	sig := sign(seed, msg)   // 64byte の署名を返す想定

	// 1. 公開鍵がバイト一致
	if !bytes.Equal(pub, stdPub) {
		t.Errorf("public key mismatch\n got=%x\nwant=%x", pub, stdPub)
	}
	// 2. 署名がバイト一致（Ed25519 は決定的なので一致するはず）
	if !bytes.Equal(sig, stdSig) {
		t.Errorf("signature mismatch\n got=%x\nwant=%x", sig, stdSig)
	}
	// 3. 自作 verify が正しい署名を受理する
	if !verify(pub, msg, sig) {
		t.Error("verify rejected a valid signature")
	}
	// 4. msg を1bit変えたら拒否する
	bad := append([]byte(nil), msg...)
	bad[0] ^= 1
	if verify(pub, bad, sig) {
		t.Error("verify accepted a tampered message")
	}
}
