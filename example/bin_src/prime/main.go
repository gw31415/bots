package main

import (
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

const (
	test_count = 200
)

func main() {

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
				num, ok := num.SetString(num_str, 10)
				//数値にならなければ
				if !ok {
					result = append(result,
						&pb.OutputMedia{
							Type:  pb.OutputMedia_UTF8,
							Data:  []byte("parse error."),
							Error: 1,
						},
					)
				} else {
					//素数判定
					if num.ProbablyPrime(test_count) {
						result = append(result,
							&pb.OutputMedia{
								Type: pb.OutputMedia_UTF8,
								Data: []byte(num_str + " is prime."),
							},
						)
					} else {
						result = append(result,
							&pb.OutputMedia{
								Type: pb.OutputMedia_UTF8,
								Data: []byte(num_str + " is NOT prime."),
							},
						)
					}
				}
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
