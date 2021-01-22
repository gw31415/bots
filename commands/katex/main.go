package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/proto"
	bots_pb "github.com/gw31415/bots/proto"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			var errmsg string
			switch err := err.(type) {
			case error:
				errmsg = err.Error()
			case string:
				errmsg = err
			default:
				errmsg = fmt.Sprint(err)
			}
			out_pb := &bots_pb.Output{
				Msgs: []*bots_pb.BotMsg{
					{
						Medias: []*bots_pb.OutputMedia{
							{
								Type:  bots_pb.OutputMedia_UTF8,
								Data:  []byte(errmsg),
								Error: 1,
							},
						},
					},
				},
			}
			//シリアライズしてバイト列にする
			out, _ := proto.Marshal(out_pb)
			os.Stdout.Write(out)
		}
	}()

	if len(os.Args) > 1 {
		var help bots_pb.Help
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			help = bots_pb.Help{
				Usage:            "[KaTeX]",
				ShortDescription: "Render KaTeX",
				LongDescription:  "Render KaTeX.",
				SourceUrl:        "https://github.com/gw31415/bots/tree/master/commands/katex",
			}
		}

		//シリアライズしてバイト列にする
		out, err := proto.Marshal(&help)
		if err != nil {
			panic(err)
		}

		//出力
		os.Stdout.Write(out)
		return
	}

	//標準入力から全てバイト列で読み込む
	arg, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	//バイト列をデシリアライズして扱えるようにする
	in_pb := &bots_pb.Input{}
	if err := proto.Unmarshal(arg, in_pb); err != nil {
		panic(err)
	}

	msgs := []*bots_pb.BotMsg{}

	//文字列データを抽出して追加
	for _, m := range in_pb.Media {
		if m.Type == bots_pb.InputMedia_UTF8 {
			KateX := strings.TrimSpace(*(*string)(unsafe.Pointer(&m.Data)))
			png, err := katex(KateX)
			if err != nil {
				panic(err)
			}
			defer os.Remove(png.Name())
			buf, err := ioutil.ReadAll(png)
			msgs = append(msgs,
				&bots_pb.BotMsg{
					Medias: []*bots_pb.OutputMedia{
						{
							Data:     buf,
							Type:     bots_pb.OutputMedia_FILE,
							Filename: filepath.Base(png.Name()),
						},
					},
				},
			)
		}
	}

	//出力データとしてまとめる
	out_pb := &bots_pb.Output{
		Msgs: msgs,
	}

	//シリアライズしてバイト列にする
	out, err := proto.Marshal(out_pb)
	if err != nil {
		panic(err)
	}

	//出力
	os.Stdout.Write(out)
}
