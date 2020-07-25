package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

/*
bots/
 ├-bin/          (ここに各コマンドに対応する実行ファイルが格納される)
 ├-proto/        (protobufの定義ファイルたち)
 ├-for_discord*  (各実装)
 └-for_line*     (各実装)

各実装とコマンドは標準入出力でやりとりする
コマンドへの入力は文字列(コマンドの引数にあたる)
コマンドからの出力はprotobufによってシリアライズされたバイト列

*/

const bindir_name = "bin"

var pb2json = jsonpb.Marshaler{Indent: "    "}

func main() {
	if len(os.Args) != 2 {
		panic("cmd name is required.")
	}
	if !regexp.MustCompile(`\w+`).Match([]byte(os.Args[1])) {
		panic("\\w+\n")
	}
	wd, err := os.Getwd()
	if err != nil {
		panic("unknown error.\n")
	}

	//呼びだすコマンドを設定
	//同一ディレクトリのbin内の実行ファイルからCmd構造体を作成
	cmd := exec.Command(filepath.Join(wd, bindir_name, os.Args[1]))

	//シリアライズされた出力はここに流しこむ
	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	//エラー出力はこっち
	stderr := bytes.NewBuffer([]byte{})
	cmd.Stderr = stderr

	//入力文字列は, 今回はシェルの標準入力から流してくる
	cmd.Stdin = os.Stdin

	//実行&待機
	err = cmd.Run()

	//コマンド側でエラーがおきたときはこのブロック
	if err != nil {
		fmt.Println(err)
		os.Stderr.Write(stderr.Bytes())
		os.Exit(1)
	}

	//ここで宣言したmsg_pbにバイト配列から変換された出力がぶっこまれる
	msg_pb := &pb.Output{}
	if err := proto.Unmarshal(out.Bytes(), msg_pb); err != nil {
		os.Exit(1)
	}
	//今回はprotobufをjsonにして標準出力に表示
	pb2json.Marshal(os.Stdout, msg_pb)
}
