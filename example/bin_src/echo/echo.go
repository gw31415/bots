package main

import (
	"bufio"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

var stdin = bufio.NewScanner(os.Stdin)

func main() {
	arg := ""
	for stdin.Scan() {
		arg += stdin.Text() + "\n"
	}
	if stdin.Err() != nil {
		// non-EOF error.
		panic(stdin.Err())
	}
	msg_pb := &pb.Output{
		Msgs: []*pb.BotMsg{
			{
				Media: &pb.BotMsg_Out{
					Out: arg,
				},
			},
		},
	}
	out, err := proto.Marshal(msg_pb)
	if err != nil {
		panic(err)
	}
	os.Stdout.Write(out)
}
