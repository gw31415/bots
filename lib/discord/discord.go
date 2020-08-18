/*
Discord用のbotsクライアントです.
*/
package discord

import (
	"fmt"
	"strings"
	"unicode"
	"unsafe"

	"github.com/bwmarrin/discordgo"
	"github.com/gw31415/bots/lib"
	"github.com/gw31415/bots/proto"
)

// Discordの一つのセッションに対処するインスタンスです.
type Instance struct {
	discord *discordgo.Session
}

// tokenを指定し, 新しいセッションを開始します.
func NewInstance(token string) (*Instance, error) {
	dc, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Instance{
		discord: dc,
	}, nil
}

// テキストチャットでメッセージを受けとったときに実行される関数です.
// matchは対処を必要とするかどうかを表します.
// cmdは実行されるコマンド名を表します.
// argは実行されるコマンドに引き渡すInput構造体を表します.
type MessagePreprocessor func(msg string) (match bool, cmd string, arg *proto.Input)

// 特定のPrefixに反応するMessagePreprocessorを返します.
func Preprocessor_prefix(prefix string) MessagePreprocessor {
	return func(msg string) (match bool, cmd string, arg *proto.Input) {
		if !strings.HasPrefix(msg, prefix) {
			return false, "", nil
		}
		msg = msg[len(prefix):]
		if unicode.IsSpace([]rune(msg)[0]) {
			return false, "", nil
		}
		match = true
		cmd = strings.Fields(msg)[0]
		arg = &proto.Input{
			Media: []*proto.InputMedia{
				{
					Data: []byte(msg[len(cmd):]),
					Type: proto.InputMedia_UTF8,
				},
			},
			Prefix: prefix,
		}
		return
	}
}

// BotMsgとmessageのデータからmessageを作る関数を表します.
type MessageStyle func(pb_out *proto.BotMsg, message_data *discordgo.MessageCreate) (embed []string)

// BotMsgとmessageのデータからembed_messageを作る関数を表します.
type MessageEmbedStyle func(pb_out *proto.BotMsg, message_data *discordgo.MessageCreate) (embed []*discordgo.MessageEmbed)

// 1メッセージ1レスポンスのサービスを表します
type UnaryServiceConfig struct {
	preprocessor MessagePreprocessor
	eStyle       MessageEmbedStyle
	style        MessageStyle
}

// Embedメッセージを返すようなUnaryServiceを返します.
func NewEmbedUnaryServiceConfig(preprocessor MessagePreprocessor, style MessageEmbedStyle) (config UnaryServiceConfig) {
	config.preprocessor = preprocessor
	config.eStyle = style
	return
}

// テキストメッセージを返すようなUnaryServiceを返します.
func NewUnaryServiceConfig(preprocessor MessagePreprocessor, style MessageStyle) (config UnaryServiceConfig) {
	config.preprocessor = preprocessor
	config.style = style
	return
}

// Instanceに新しいUnaryServiceを設定します.
func (instance *Instance) SetUnaryService(bots *lib.BotsHandler, config UnaryServiceConfig) {
	handler := func(session *discordgo.Session, message *discordgo.MessageCreate) {
		if message.Author.Bot || message.Author.ID == session.State.User.ID { // 自身とBotのメッセージは無視
			return
		}
		match, cmd, arg := config.preprocessor(message.Content)
		if match {
			botscmd, err := bots.GetCommand(cmd)
			if err != nil { // 見つからない, または壊れている
				return
			}
			out_pb, err := botscmd.Run(arg)
			if err != nil { // 起動に失敗, または壊れている
				return
			}
			if config.style == nil { // Embedメッセージを返すとき
				for _, msg := range out_pb.Msgs {
					msg_list := config.eStyle(msg, message)
					for _, m := range msg_list {
						session.ChannelMessageSendEmbed(message.ChannelID, m)
					}
				}
			} else { // テキストメッセージを返すとき
				for _, msg := range out_pb.Msgs {
					msg_list := config.style(msg, message)
					for _, m := range msg_list {
						session.ChannelMessageSend(message.ChannelID, m)
					}
				}
			}
		}
	}
	instance.discord.AddHandler(handler)
	instance.discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
}

// Instanceを起動します.
func (instance *Instance) Open() error {
	return instance.discord.Open()
}

// Instanceを停止します.
func (instance *Instance) Close() error {
	return instance.discord.Close()
}

// デフォルトのMessageEmbedStyleです.
func MessageEmbedStyle_default(msg *proto.BotMsg, message *discordgo.MessageCreate) (embeds []*discordgo.MessageEmbed) {
	error_count := 0
	embed := &discordgo.MessageEmbed{}
	field := &discordgo.MessageEmbedField{}
	for i, media := range msg.Medias {
		if media.Error != 0 {
			error_count++
		}
		switch media.Type {
		case proto.OutputMedia_UTF8:
			str := *(*string)(unsafe.Pointer(&media.Data))
			if media.Level == 0 {
				line := ""
				if field.Value != "" {
					line = "\n"
				}
				line += str
				if media.LongCode {
					line = fmt.Sprintf("```\n%s\n```", strings.ReplaceAll(line, "```", " `` "))
				} else if media.ShortCode {
					qs := []rune{}
					for i := 0; i < count_bquates(line)+1; i++ {
						qs = append(qs, '`')
					}
					bf := line
					line = string(qs)
					if bf[0] == '`' {
						line += " "
					}
					line += bf
					if bf[len(bf)-1] == '`' {
						line += " "
					}
					line += string(qs)
				} else if media.Spoiled {
					line = "||" + strings.ReplaceAll(line, "||", "\\|\\|")
				}
				field.Value += line
			} else {
				field.Name = str
			}
		}
		if i+1 == len(msg.Medias) || !media.ExtendField {
			embed.Fields = append(embed.Fields, field)
			field = &discordgo.MessageEmbedField{}
		}
		if i+1 == len(msg.Medias) || media.Type == proto.OutputMedia_FILE {
			if len(embed.Fields) == 1 && embed.Fields[0].Name == "" {
				embed.Description = embed.Fields[0].Value
				embed.Fields = []*discordgo.MessageEmbedField{}
			} else {
				for _, f := range embed.Fields {
					if f.Name == "" {
						f.Name = "(out)"
					}
				}
			}
			embed.Color = int(msg.Color & 0x00ffffff)
			author := &discordgo.MessageEmbedAuthor{}
			author.Name = message.Message.Author.Username
			author.IconURL = message.Author.AvatarURL("")
			embed.Author = author
			if media.Type == proto.OutputMedia_FILE {
			} else {
				embeds = append(embeds, embed)
				embed = &discordgo.MessageEmbed{}
			}
		}
	}
	return embeds
}

func count_bquates(st string) int {
	mx := 0
	i := 0
	for _, x := range st {
		if x == '`' {
			i++
			if mx < i {
				mx = i
			}
		} else {
			i = 0
		}
	}
	return mx
}
