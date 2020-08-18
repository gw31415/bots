/*
botsの仕様にあわせたライブラリです
*/
package lib

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"unsafe"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
)

var (
	commandBrokenError   = errors.New("command is broken.")
	commandNotFoundError = errors.New("command not found.")
	commandStartingError = errors.New("command failed to start.")
)

// コマンドのハンドラ
type BotsHandler struct {
	bin_dir string
}

// bin_dirにコマンドが配置されたディレクトリのパスを指定, ハンドラを返します
func NewBotsHandler(bin_dir string) (*BotsHandler, error) {
	if f, err := os.Stat(bin_dir); os.IsNotExist(err) || !f.IsDir() {
		return nil, errors.New("direcrory does not exist.")
	}
	p, err := filepath.Abs(bin_dir)
	if err != nil {
		return nil, err
	}
	return &BotsHandler{
		bin_dir: p,
	}, nil
}

// Commandはコマンドを表します
type Command struct {
	ex   *exec.Cmd
	help *pb.Help
}

// コマンドのオブジェクトを返します
func (handler *BotsHandler) GetCommand(cmd_name string) (*Command, error) {
	cmd_path := fmt.Sprintf("%s/%s", handler.bin_dir, cmd_name)
	f, err := os.Stat(cmd_path)
	if os.IsNotExist(err) {
		return nil, commandNotFoundError
	}
	if f.IsDir() {
		cmd_path += "/" + cmd_name
		if _, err = os.Stat(cmd_path); os.IsNotExist(err) {
			return nil, commandNotFoundError
		}
	}
	help_getter := exec.Command(cmd_path, "-h")
	help_bin := bytes.NewBuffer([]byte{})
	help_getter.Stdout = help_bin
	if help_getter.Run() != nil {
		return nil, commandBrokenError
	}
	help_pb := &pb.Help{}
	if err := proto.Unmarshal(help_bin.Bytes(), help_pb); err != nil {
		return nil, commandBrokenError
	}
	return &Command{
		ex:   exec.Command(cmd_path),
		help: help_pb,
	}, nil
}

// コマンドを起動します
func (cmd *Command) Run(in_pb *pb.Input) (*pb.Output, error) {
	stdin, err := cmd.ex.StdinPipe()
	out := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})
	cmd.ex.Stdout = out
	cmd.ex.Stderr = stderr
	if err != nil {
		return nil, err
	}
	if cmd.ex.Start() != nil {
		return nil, commandStartingError
	}
	data, err := proto.Marshal(in_pb)
	if err != nil {
		return nil, commandBrokenError
	}
	stdin.Write(data)
	stdin.Close()
	if err := cmd.ex.Wait(); err != nil {
		bye := stderr.Bytes()
		return nil, errors.New(fmt.Sprintf("%s\n%s", *(*string)(unsafe.Pointer(&bye)), err.Error()))
	}
	msg_pb := &pb.Output{}
	if err := proto.Unmarshal(out.Bytes(), msg_pb); err != nil {
		return nil, commandBrokenError
	}
	return msg_pb, nil
}

// Helpを返します
func (cmd *Command) GetHelp() *pb.Help {
	return cmd.help
}
