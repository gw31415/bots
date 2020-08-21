package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

func Gcd(nums ...*big.Int) *big.Int {
	switch len(nums) {
	case 2:
		if nums[0].Cmp(nums[1]) == -1 {
			n := nums[1]
			nums[1] = nums[0]
			nums[0] = n
		}
		if nums[1].Sign() == 0 {
			return nums[0]
		}
		return Gcd(nums[1], nums[0].Mod(nums[0], nums[1]))
	case 1:
		return nums[0]
	default:
		g := Gcd(nums[0], nums[1])
		nums = nums[1:]
		nums[0] = g
		return Gcd(nums...)
	}
}

func SetString(n *big.Int, str string) (nw *big.Int) {
	nw, ok := n.SetString(str, 10)
	ok = ok && nw.Sign() == 1
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
				ShortDescription: "最大公約数",
				LongDescription:  "与えられた正の整数の最大公約数を返します.",
				SourceUrl:        "https://github.com/gw31415/bots/tree/master/commands/gcd",
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

	var medias []*pb.OutputMedia
	//メディア配列を回す
	for _, m := range in_pb.Media {
		if m.Type == pb.InputMedia_UTF8 {
			strs := strings.Fields(string(m.Data))
			result := &big.Int{}
			SetString(result, strs[0])
			for _, num_str := range strs {
				//数値にする
				num := new(big.Int)
				//数値にならなければ
				SetString(num, num_str)
				result = Gcd(result, num)
			}
			medias = append(medias, &pb.OutputMedia{
				Type: pb.OutputMedia_UTF8,
				Data: []byte(result.String()),
			})
		}
	}
	//出力データとしてまとめる
	out_pb := &pb.Output{
		Msgs: []*pb.BotMsg{
			{
				Medias: medias,
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
