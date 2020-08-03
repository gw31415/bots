package main

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

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
	arg, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	//文字列でメッセージをつくる
	in_pb := &pb.Input{
		Media: []*pb.InputMedia{
			{
				Data: arg,
				Type: pb.InputMedia_UTF8,
			},
		},
	}

	//メッセージのシリアライズ
	data, err := proto.Marshal(in_pb)
	if err != nil {
		panic(err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}
	//実行&待機
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	stdin.Write(data)
	stdin.Close()
	cmd.Wait()

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
