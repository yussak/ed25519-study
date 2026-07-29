# ed25519-jisaku

Ed25519 (RFC 8032) を Go でゼロから自作して理解する学習プロジェクト。
「作って学ぶ」方針で、TDD で漸進的に進める。

## これは何？
- 楕円曲線署名 Ed25519 を、標準ライブラリの `crypto/ed25519` に頼らず
  `math/big` と `crypto/sha512` だけで実装してみる。
- 目的は理解。**本番用途には使わない**（タイミング攻撃対策などは学習範囲外）。

## 進め方（TDD の最小ループ）
1関数ずつ **Red → Green → Refactor**：
1. 🔴 小さな失敗テストを1本書く
2. 🟢 通る最小限のコードを書く
3. ♻️ 整える

実装順序は **field（有限体）→ curve（曲線）→ sign（署名）**。
詳しい方針は [`CLAUDE.md`](./CLAUDE.md)、やることリストは
[`ed25519-study/TODO.md`](./ed25519-study/TODO.md)。

## 動かし方
```sh
cd ed25519-study
go run .          # main が動く
go test ./...     # ここで TDD を回す
```

## 構成
```
ed25519-study/
├── field.go   有限体 GF(2^255-19)
├── curve.go   ねじれ Edwards 曲線と群演算
├── sign.go    鍵生成 / 署名 / 検証 (RFC 8032 §5.1)
├── main.go    デモ用エントリポイント
└── *_test.go  各段のテスト
```

## 参考
- RFC 8032 (EdDSA), §5.1 が Ed25519
