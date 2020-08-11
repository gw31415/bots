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
	"io/ioutil"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/gw31415/bots/proto"
	"github.com/spf13/cobra"
)

// inbinCmd represents the inbin command
var inbinCmd = &cobra.Command{
	Use:   "inbin",
	Short: "Output cmd_in binary",
	Long: `Make and output cmd_in binary.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		//入力文字列は, 今回はシェルの標準入力から流してくる
		arg, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
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
		data, err := proto.Marshal(in_pb)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(data)
		return err
	},
}

func init() {
	rootCmd.AddCommand(inbinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// inbinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// inbinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
