package main

import (
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
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

	//出力データをまとめる
	out_pb := &pb.Output{
		Msgs: []*pb.BotMsg{
			{
				Medias: []*pb.OutputMedia{
					{
						Data: in_pb.Media[0].Data,
						Type: pb.OutputMedia_UTF8,
					},
				},
			},
		},
	}

	//シリアライズしてバイト列にする
	out, err := proto.Marshal(out_pb)
	if err != nil {
		panic(err)
	}

	//出力
	os.Stdout.Write(out)

}
