package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

func main() {

	if len(os.Args) > 1 {
		var help pb.Help
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			help = pb.Help{
				Usage:            "(CMD)",
				ShortDescription: "Show help",
				LongDescription:  "Show help",
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

	medias := []*pb.OutputMedia{}

	//文字列データを抽出して追加
	for _, m := range in_pb.Media {
		if m.Type == pb.InputMedia_UTF8 {
			for _, cmd := range strings.Fields(*(*string)(unsafe.Pointer(&m.Data))) {
				wd, err := os.Executable()
				wd = filepath.Dir(wd)
				if err != nil {
					panic(err)
				}
				cmd_path := fmt.Sprintf("%s/%s", wd, cmd)
				botcmd := exec.Command(cmd_path, "-h")
				out := bytes.NewBuffer([]byte{})
				botcmd.Stdout = out
				err = botcmd.Run()
				if err != nil {
					panic(err)
				}
				help_pb := &pb.Help{}
				if err := proto.Unmarshal(out.Bytes(), help_pb); err != nil {
					panic(err)
				}
				medias = append(medias,
					&pb.OutputMedia{
						Data:        []byte(cmd),
						Type:        pb.OutputMedia_UTF8,
						Level:       1,
						ExtendField: true,
					},
					&pb.OutputMedia{
						Data:        []byte(fmt.Sprintf("%s%s %s",in_pb.Prefix, cmd, help_pb.Usage)),
						Type:        pb.OutputMedia_UTF8,
						ExtendField: true,
						ShortCode: true,
					},
					&pb.OutputMedia{
						Data: []byte(help_pb.LongDescription),
						Type: pb.OutputMedia_UTF8,
					},
				)
			}

		} else {
			medias = append(medias,
				&pb.OutputMedia{
					Data:        []byte("Error"),
					Type:        pb.OutputMedia_UTF8,
					Error:       1,
					Level:       1,
					ExtendField: true,
				},
				&pb.OutputMedia{
					Data:  []byte("parse error"),
					Type:  pb.OutputMedia_UTF8,
					Error: 1,
				},
			)
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
