package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"github.com/mathlava/bigc"
	"github.com/mathlava/bigc/math/rho"
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
			out_pb := &pb.Output{
				Msgs: []*pb.BotMsg{
					{
						Medias: []*pb.OutputMedia{
							{
								Type:  pb.OutputMedia_UTF8,
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
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			help := pb.Help{
				Usage:            "[Numbers]",
				ShortDescription: "素微分します",
				LongDescription:  "素微分をします.",
				SourceUrl:        "https://github.com/gw31415/bots/tree/master/commands/ad",
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

	//メディア配列を回す
	for _, m := range in_pb.Media {
		if m.Type == pb.InputMedia_UTF8 {
			var result []*pb.OutputMedia
			for _, num_str := range strings.Fields(string(m.Data)) {
				//数値にする
				num, err := bigc.ParseString(num_str)
				if err != nil {
					panic("parse error.")
				}
				rho.ArithmeticDerivative(num)
				result = append(result,
					&pb.OutputMedia{
						Type:        pb.OutputMedia_UTF8,
						Data:        []byte(num_str),
						ExtendField: true,
						Level:       1,
					},
					&pb.OutputMedia{
						Type: pb.OutputMedia_UTF8,
						Data: []byte(num.String()),
					},
				)
			}
			msgs = append(msgs,
				&pb.BotMsg{
					Medias: result,
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
