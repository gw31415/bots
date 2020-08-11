/*
Copyright © 2020 Amadeus_vn <git@amas.dev>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"

	"github.com/spf13/cobra"
)

const bindir_name = "bin"

var (
	pb2json  = jsonpb.Marshaler{Indent: "	"}
	json_out = false
	json_in  = false
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test bots commands",
	Long:  `Test commands of bots by sending input and reseive output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("one cmd name is required.")
		}
		if !regexp.MustCompile(`\w+`).Match([]byte(args[0])) {
			return errors.New("\\w+\n")
		}
		wd, err := os.Getwd()
		if err != nil {
			return errors.New("unknown error.\n")
		}

		//呼びだすコマンドを設定
		//同一ディレクトリのbin内の実行ファイルからCmd構造体を作成
		botcmd := exec.Command(filepath.Join(wd, bindir_name, args[0]))

		//シリアライズされた出力はここに流しこむ
		out := bytes.NewBuffer([]byte{})
		botcmd.Stdout = out
		//エラー出力はこっち
		stderr := bytes.NewBuffer([]byte{})
		botcmd.Stderr = stderr

		//入力文字列は, 今回はシェルの標準入力から流してくる
		arg, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}

		var data []byte
		if json_in {
			return errors.New("json input is not supported yet.")
		} else {

			//文字列でメッセージをつくる
			in_pb := &pb.Input{
				Prefix: "[PREFIX]",
				Media: []*pb.InputMedia{
					{
						Data: arg,
						Type: pb.InputMedia_UTF8,
					},
				},
			}

			//メッセージのシリアライズ
			data, err = proto.Marshal(in_pb)
			if err != nil {
				return err
			}
		}

		stdin, err := botcmd.StdinPipe()
		if err != nil {
			return err
		}

		//実行&待機
		err = botcmd.Start()
		if err != nil {
			return err
		}

		stdin.Write(data)
		stdin.Close()
		botcmd.Wait()

		//コマンド側でエラーがおきたときはこのブロック
		if err != nil {
			fmt.Println(err)
			os.Stderr.Write(stderr.Bytes())
			os.Exit(1)
		}

		//ここで宣言したmsg_pbにバイト配列から変換された出力がぶっこまれる
		msg_pb := &pb.Output{}
		if err := proto.Unmarshal(out.Bytes(), msg_pb); err != nil {
			os.Exit(1)
		}

		if json_out {
			//protobufをjsonにして標準出力に表示
			pb2json.Marshal(os.Stdout, msg_pb)
		} else {
			fmt.Println()
			for _, msg := range msg_pb.Msgs {
				for _, media := range msg.Medias {
					switch media.Type {
					case pb.OutputMedia_UTF8:
						fmt.Printf("UTF8: `%s`", string(media.Data))
						fmt.Println()
					case pb.OutputMedia_FILE:
						if media.Filename == "" {
							fmt.Println("File: BROKEN")
							break
						}
						if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
							os.Mkdir("./tmp", 0755)
						}
						if err != nil {
							fmt.Printf("File: (error) %s", err.Error())
							fmt.Println()
							break
						}
						file, err := os.OpenFile("./tmp/"+media.Filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
						if err != nil {
							fmt.Printf("File: (error) %s", err.Error())
							fmt.Println()
							break
						}
						if _, e := file.Write(media.Data); e != nil {
							fmt.Printf("File: (error) %s", e.Error())
						}
						fmt.Printf("File: ./fmt/%s", media.Filename)
						fmt.Println()
					}
				}
				fmt.Println()
			}
		}
		return nil

	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolVarP(&json_in, "json-in", "i", false, "Use json input")
	testCmd.Flags().BoolVarP(&json_out, "json-out", "o", false, "Use json input")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
