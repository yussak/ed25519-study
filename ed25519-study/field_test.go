package main

import "testing"

// field_test.go : 有限体の演算をここで検証する。
//
// いまは「go test が通る最小状態」を保つためのプレースホルダのみ。
// field.go の実装を進めたら、下の t.Skip を消して本テストを書く。

func TestFieldPlaceholder(t *testing.T) {
	// TODO: feMul(a, feInv(a)) == 1 を確認するテストに置き換える。
	//   例) a を適当に取り、new(big.Int) で 1 と等しいか比較する。
	t.Skip("field.go 未実装: 実装したらこの Skip を消して本テストを書く")
}
