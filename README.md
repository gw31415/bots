# bots
## コンセプト
* プラットフォームによらないコマンドハンドラ
* 開発言語によらず開発できるコマンド群
* コマンドハンドラ毎に設定できる出力デザイン
* 集団および分散開発のしやすさ
* 他者の開発したコマンドの導入のしやすさ

## 仕様
### コマンドハンドラ
1. あるただ一つの特定のディレクトリを読みこむ. それを読みこみディレクトリとする.
2. 読みこみディレクトリ直下に読みこまれる全てのコマンドが置かれる.
3. 読みこみディレクトリにはコマンドの他のファイルを配置しない.
4. ある特定のコマンドを呼びだす操作は, 引数なしで呼びだすコマンドを実行することによって行なわれる.
5. ある特定のコマンドのヘルプを呼びだす操作は, ただひとつの引数 `--help` または `-h` で呼びだすコマンドを実行することによって行なわれる.
6. コマンドのヘルプは, `proto/help.proto` に基づきシリアライズされている.
7. コマンドへの入力は, コマンドの実行中に標準入力に書きこむ.
8. コマンドの入力は `proto/cmdin.proto` に基づきシリアライズされている.
9. コマンドの出力は, コマンドの実行中に標準出力から読みこむ.
10. コマンドの出力は `proto/cmdout.proto` に基づきシリアライズされている.

### コマンド
1. 実行ファイルとする.
2. コマンドのファイル名はコマンド名とする.
3. コマンド名の全体は `\w+` に適合する.
4. コマンドの通常動作は引数なしで呼びだされると行なわれる.
5. コマンドのヘルプを呼びだす操作は, ただひとつの引数 `--help` または `-h` で呼びだすコマンドを実行することによって行なわれる.
6. コマンドのヘルプは, `proto/help.proto` に基づきシリアライズされている.
7. コマンドへの入力は, コマンドの実行中に標準入力に書きこむ.
8. コマンドの入力は `proto/cmdin.proto` に基づきシリアライズされている.
9. コマンドの出力は, コマンドの実行中に標準出力から読みこむ.
10. コマンドの出力は `proto/cmdout.proto` に基づきシリアライズされている.
