package main

import (
	"io/ioutil"
	"os"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"github.com/gw31415/hubcontrib"
)

const (
	svg_filename = "hubcontrib.svg"
)

func main() {

	if len(os.Args) > 1 {
		var help pb.Help
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			help = pb.Help{
				Usage:            "[ID]",
				ShortDescription: "Show the GitHub contribution graph",
				LongDescription:  "Show year-long GitHub contribution graph.",
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
	in_pb := &pb.Input{}
	if err := proto.Unmarshal(arg, in_pb); err != nil {
		panic(err)
	}

	msgs := []*pb.BotMsg{}

	//文字列データを抽出して追加
	for _, m := range in_pb.Media {
		if m.Type == pb.InputMedia_UTF8 {
			uname := strings.TrimSpace(*(*string)(unsafe.Pointer(&m.Data)))
			svg, err := hubcontrib.Svg(uname)
			if err != nil {
				msgs = append(msgs,
					&pb.BotMsg{
						Medias: []*pb.OutputMedia{
							{
								Data:  []byte(err.Error()),
								Type:  pb.OutputMedia_UTF8,
								Error: 1,
							},
						},
					},
				)
			}
			msgs = append(msgs,
				&pb.BotMsg{
					Medias: []*pb.OutputMedia{
						{
							Data: []byte(svg),
							Type: pb.OutputMedia_FILE,
							Filename: svg_filename,
						},
					},
				},
			)
		}
	}

	//出力データとしてまとめる
	out_pb := &pb.Output{
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
