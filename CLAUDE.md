# CLAUDE.md

TDD で対応

まずは main.go に全部まとめて実装する
分割（field / curve / sign）は中身が分かってきてから後でやる
理由: 複雑さを減らすため、まず1ファイルで動かして構造が見えてから分ける

当面のゴールは main_test.go の TestAgainstStdEd25519 を通すこと（Ed25519 コア完成の定義）
