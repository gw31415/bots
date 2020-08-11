# botcli
botsコマンドをテストするコマンドラインツール。

## 準備
botsのプロトコルに準拠した実行ファイル(もしくはリンク)を同一ディレクトリに作成した`./bin`下に配置する。

## サブコマンド

| サブコマンド | 使い方 |
:---: | :---
test | `botcli test [コマンド名]`として`UTF8`でメッセージを作成する標準入力に入る
help | `botcli help [サブコマンド名]`としてサブコマンドの使い方を表示する