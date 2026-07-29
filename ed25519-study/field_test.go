package main

import (
	"math/big"
	"testing"
)

// field_test.go : 有限体の演算をここで検証する。
//
// 【TDD 最初の1歩：worked example】
// 下の TestFeAdd は「お手本」の1本。まだ field.go に p も feAdd も無いので、
//   go test  →  RED（コンパイルエラー）になるはず。まずそれを確認する。
// 次に field.go に p と feAdd を書いて GREEN にする。
// feSub / feMul / feInv は、このパターンを真似て自分で書く（下の TODO）。

func TestFeAdd(t *testing.T) {
	// ケース1: 素直な足し算（まだ mod は跨がない）
	got := feAdd(big.NewInt(2), big.NewInt(3))
	if got.Cmp(big.NewInt(5)) != 0 {
		t.Errorf("feAdd(2,3) = %v, want 5", got)
	}

	// ケース2: p を跨ぐ足し算。 (p-1) + 3 = p + 2 ≡ 2 (mod p)
	//   p の実際の桁を知らなくても、この等式で正しさを確認できる。
	pm1 := new(big.Int).Sub(p, big.NewInt(1)) // p-1
	got = feAdd(pm1, big.NewInt(3))
	if got.Cmp(big.NewInt(2)) != 0 {
		t.Errorf("feAdd(p-1,3) = %v, want 2", got)
	}
}

// TODO: TestFeSub を自分で書く。ヒント: 3 - 5 は 0 未満。mod p で正の値に畳めているか？
// TODO: TestFeMul を自分で書く。ヒント: (p-1) * 2 ≡ p-2 (mod p) など。
// TODO: TestFeInv を自分で書く。property test: feMul(a, feInv(a)) == 1。
