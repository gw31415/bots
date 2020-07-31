package main

import (
	"bufio"
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

var stdin = bufio.NewScanner(os.Stdin)

func main() {
	arg, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	in_pb := &pb.Input{}
	if err := proto.Unmarshal(arg, in_pb); err != nil {
		panic(err)
	}
	out_pb := &pb.Output{
		Msgs: []*pb.BotMsg{
			{
				Medias: []*pb.OutputMedia{
					{
						Data: arg,
						Type: pb.OutputMedia_UTF8,
					},
				},
			},
		},
	}
	out, err := proto.Marshal(out_pb)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
}
