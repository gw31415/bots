package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"github.com/mathlava/bigc/math/rho"
	"github.com/mathlava/bigc"
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
				Ad(num)
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

var (
	val1 = big.NewInt(1)
	val0 = big.NewInt(0)
	val1_2 = big.NewRat(1,2)
)

func Ad(num *bigc.BigC) *bigc.BigC {
	if num.IsReal() {
		ad_rat(num.Real())
		return num
	}
	if num.IsPureImag() {
		ad_rat(num.Imag())
		return num
	}
	cache := new(bigc.BigC).Set(num).AbsSq()
	num.Imag().Quo(num.Imag(),cache).Mul(num.Imag(), val1_2)
	num.Real().Quo(num.Real(),cache).Mul(num.Real(), val1_2)
	ad_rat(cache)
	num.Real().Mul(num.Real(), cache)
	num.Imag().Mul(num.Imag(), cache)
	return num
}

func ad_rat(num *big.Rat) *big.Rat {
	sign := num.Sign()
	num.Abs(num)
	ba_ := ad_int(new(big.Int).Set(num.Denom()))
	ba_.Mul(ba_, num.Num())
	ad_int(num.Num()).Mul(num.Num(), num.Denom()).Sub(num.Num(), ba_)
	num.Denom().Mul(num.Denom(), num.Denom())
	if sign == -1 {
		num.Neg(num)
	}
	return num
}

func ad_int(num *big.Int) *big.Int {
	sign := num.Sign()
	num.Abs(num)
	if num.Cmp(val1) == 0 {
		return num.Set(val0)
	}
	if sign == 0 {
		panic("undefined.")
	}
	cache := new(big.Int).Set(num)
	num.Set(val0)
	add := new(big.Int)
	for _, p := range rho.Primes(new(big.Int).Set(cache)) {
		add.Set(cache)
		add.Div(add, p)
		num.Add(num, add)
	}
	if sign == -1 {
		num.Neg(num)
	}
	return num
}
