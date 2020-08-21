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
		var help pb.Help
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			help = pb.Help{
				Usage:            "[CMD]",
				ShortDescription: "Source url",
				LongDescription:  "ソースコードのURLを返します.",
				SourceUrl:        "https://github.com/gw31415/bots/tree/master/commands/source",
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

	wd, err := os.Executable()
	if err != nil {
		panic(err)
	}
	wd = filepath.Dir(wd)
	cmd_names := []string{}
	for _, m := range in_pb.Media {
		if m.Type == pb.InputMedia_UTF8 {
			for _, cmd := range strings.Fields(*(*string)(unsafe.Pointer(&m.Data))) {
				cmd_names = append(cmd_names, cmd)
			}
		} else {
			cmd_names = append(cmd_names, "")
		}
	}
	if len(cmd_names) == 0 {
		panic("please specify a command.")
	}

	for _, cmd := range cmd_names {
		if cmd == "" {
			panic("parse error.")
		}
		cmd_path := fmt.Sprintf("%s/%s", wd, cmd)
		fInfo, _ := os.Stat(cmd_path)
		if fInfo.IsDir() {
			cmd_path += "/" + cmd
		}
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
		if help_pb.SourceUrl == "" {
			panic("source url not found.")
		}
		medias = append(medias,
			&pb.OutputMedia{
				Data:        []byte(cmd),
				Type:        pb.OutputMedia_UTF8,
				Level:       1,
				ExtendField: true,
			},
			&pb.OutputMedia{
				Data: []byte("Source Url: " + help_pb.SourceUrl),
				Type: pb.OutputMedia_UTF8,
			},
		)

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
