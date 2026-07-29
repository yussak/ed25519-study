---
name: grill-me
description: Interrogate the user about a plan or design until full mutual understanding is reached. Use when the user wants their plan, design, or approach pressure-tested through rigorous one-at-a-time questioning before implementation — triggered by "grill me", "質問攻めにして", "設計を詰めたい", or asking to be challenged on a plan.
---

# grill-me

計画・設計の**あらゆる側面**について、共通認識に達するまで徹底的にユーザーに質問する。

## 進め方

1. **設計のツリーを枝分かれの先まで一つひとつたどる。** 決定事項間の依存関係を順番に解決していく。上流の決定（前提・目的・制約）から下流（具体的な実装選択）へ。ある決定が別の決定に依存するなら、依存先を先に片付ける。

2. **質問は一度に一つずつ。** まとめて複数投げない。ユーザーの回答を受けてから次の質問へ進む。

3. **各質問に、あなたの推奨する回答を必ず添える。** 「どう思いますか？」だけで終わらせず、「私の推奨は X。理由は〜。あなたはどうしますか？」の形にする。推奨には根拠を付ける。

4. **コードベースで答えが出る質問は、質問せず自分で調べる。** 「今どういう構成になっているか」「既存の命名/依存はどうか」のような事実確認は、ユーザーに聞く前に Read/Grep/Glob 等で調査して確定させる。ユーザーに聞くのは、ユーザーにしか決められないこと（目的・優先順位・トレードオフの取捨選択・好み）だけ。

## スタイル

- 一問一答のテンポを保つ。共通認識に達した論点は短く確認して次へ。
- 曖昧な合意で流さない。決定は具体化してから次の枝へ。
- 全ての枝を解決し終えたら、合意事項を簡潔にまとめて確認する。
