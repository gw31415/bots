package main

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/golang/protobuf/proto"
	bots_pb "github.com/gw31415/bots/proto"
	texc_pb "github.com/gw31415/texc/proto"
	"google.golang.org/grpc"
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
				Usage:            "[source]",
				ShortDescription: "TeX文書のレンダリング",
				LongDescription:  "TeX文書をレンダリングします.\n 初期状態では数式モードではないことに注意してください.",
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
			conn, err := grpc.Dial("texc.amas.dev:3475", grpc.WithInsecure())
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			client := texc_pb.NewTexcServiceClient(conn)
			stream, err := client.Sync(context.Background())
			if err != nil {
				panic(err)
			}
			tar_data := bytes.NewBuffer([]byte{})
			tar_w := tar.NewWriter(tar_data)
			buf := new(bytes.Buffer)
			size, _ := fmt.Fprintf(buf, `
			\documentclass[uplatex]{jsarticle}

			\usepackage{amsmath,amssymb}
			\usepackage{ascmac}
			\usepackage{enumerate}
			\usepackage{mathrsfs}
			\usepackage{bm}
			\usepackage{chemfig}
			\usepackage[version=3]{mhchem}
			\usepackage{tikz}
			\usepackage{pxrubrica}
			\usepackage{siunitx}
			\usepackage{circuitikz}
			\usepackage{cancel}
			\usepackage{fancybox}

			\pagestyle{empty}

			\begin{document}
				%s
			\end{document}`, (*(*string)(unsafe.Pointer(&m.Data))))
			tar_w.WriteHeader(&tar.Header{
				Name:    "bot.tex",
				Mode:    0755,
				ModTime: time.Now(),
				Size:    int64(size),
			})
			io.Copy(tar_w, buf)
			tar_w.Close()
			in_pb := new(texc_pb.Input)
			in_pb.Data = make([]byte, 0xff)
			for {
				_, err := tar_data.Read(in_pb.Data)
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}
				stream.Send(in_pb)
			}
			stream.Send(&texc_pb.Input{
				Exec: []string{"latexmk", "bot.tex"},
			})
			stream.Send(&texc_pb.Input{
				Exec: []string{"pdfcrop", "--margins", "50 50 50 50", "bot.pdf", "bot.pdf"},
			})
			stream.Send(&texc_pb.Input{
				Exec: []string{"pdftoppm", "-singlefile", "-jpeg", "-r", "400", "bot.pdf", "bot"},
			})
			stream.Send(&texc_pb.Input{
				Dl: "bot.jpg",
			})
			stream.CloseSend()
			stdout := bytes.NewBufferString("")
			tar_dl_data := bytes.NewBuffer([]byte{})
			for {
				out, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					body_bin := stdout.Bytes()
					body_str := (*(*string)(unsafe.Pointer(&body_bin)))
					if tar_dl_data.Len() == 0 {
						i := strings.Index(body_str, "!")
						if i == -1 {
							panic(body_str)
						}
						panic(body_str[i:])
					}
					panic(err)
				}
				if out.Stdout != nil {
					stdout.Write(out.Stdout)
				}
				if out.Data != nil {
					tar_dl_data.Write(out.Data)
				}
			}
			body_bin := stdout.Bytes()
			body_str := (*(*string)(unsafe.Pointer(&body_bin)))
			if tar_dl_data.Len() == 0 {
				i := strings.Index(body_str, "!")
				if i == -1 {
					panic(body_str)
				}
				panic(body_str[i:])
			}
			jpeg_data := bytes.NewBuffer([]byte{})
			tar_reader := tar.NewReader(tar_dl_data)
			for {
				_, err := tar_reader.Next()
				/*
					if err == io.EOF {
						break
					}
				*/
				if err != nil {
					break
					//ここらへんバグ
					//panic(err)
				}
				io.Copy(jpeg_data, tar_dl_data)
			}
			msgs = append(msgs,
				&bots_pb.BotMsg{
					Medias: []*bots_pb.OutputMedia{
						{
							Data:     jpeg_data.Bytes(),
							Type:     bots_pb.OutputMedia_FILE,
							Filename: "tex.jpg",
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
