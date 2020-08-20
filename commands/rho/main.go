package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"github.com/mathlava/bigc/math/rho"
)

func SetString(n *big.Int, str string) (nw *big.Int) {
	nw, ok := n.SetString(str, 10)
	if ok {
		return nw
	} else {
		panic("parse error.")
	}
}

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
				ShortDescription: "ポラード・ロー法",
				LongDescription:  "ポラード・ロー法",
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
				num := new(big.Int)
				num = SetString(num, num_str)
				//数値にならなければ
				list := rho.Primes(num)
				sort.Slice(list, func(i, j int) bool { return list[i].Cmp(list[j]) == -1 })
				buf := make([]byte, 0)
				for i := 0; i < len(list)-1; i++ {
					buf = append(buf, list[i].String()...)
					buf = append(buf, "×"...)
				}
				buf = append(buf, list[len(list)-1].String()...)
				result = append(result,
					&pb.OutputMedia{
						Type: pb.OutputMedia_UTF8,
						Data: buf,
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
