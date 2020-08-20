/*
Discord用のbotsクライアントです.
*/
package discord

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unsafe"

	"github.com/bwmarrin/discordgo"
	"github.com/gw31415/bots/lib"
	"github.com/gw31415/bots/proto"
	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
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
type MessageStyle func(pb_out *proto.BotMsg, message_data *discordgo.MessageCreate) (message []*discordgo.MessageSend)

// 1メッセージ1レスポンスのサービスを表します
type UnaryServiceConfig struct {
	preprocessor MessagePreprocessor
	style        MessageStyle
	Timeout      time.Duration
}

// メッセージを返すようなUnaryServiceを返します.
func NewUnaryServiceConfig(preprocessor MessagePreprocessor, style MessageStyle) (config UnaryServiceConfig) {
	config.preprocessor = preprocessor
	config.style = style
	config.Timeout = 5 * time.Second
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
			session.ChannelTyping(message.ChannelID)
			done := make(chan *proto.Output, 1)
			go func() {
				out_pb, err := botscmd.Run(arg)
				if err != nil {
					out_pb = &proto.Output{}
					out_pb.Msgs = []*proto.BotMsg{
						{
							Medias: []*proto.OutputMedia{
								{
									Error: 1,
									Data:  []byte(err.Error()),
									Type:  proto.OutputMedia_UTF8,
								},
							},
						},
					}
				}
				done <- out_pb
			}()
			out_pb := &proto.Output{}
			select {
			case <-time.After(config.Timeout):
				if err := botscmd.Kill(); err != nil {
					out_pb.Msgs = []*proto.BotMsg{
						{
							Medias: []*proto.OutputMedia{
								{
									Error: 1,
									Data:  []byte(err.Error()),
									Type:  proto.OutputMedia_UTF8,
								},
							},
						},
					}
				} else {
					out_pb.Msgs = []*proto.BotMsg{
						{
							Medias: []*proto.OutputMedia{
								{
									Error: 1,
									Data:  []byte("timeout."),
									Type:  proto.OutputMedia_UTF8,
								},
							},
						},
					}
				}
			case pb := <-done:
				out_pb = pb
			}
			for _, msg := range out_pb.Msgs {
				msg_list := config.style(msg, message)
				for _, m := range msg_list {
					session.ChannelMessageSendComplex(message.ChannelID, m)
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
func MessageEmbedStyle_default(msg *proto.BotMsg, message *discordgo.MessageCreate) (sends []*discordgo.MessageSend) {
	// ランタイムエラーの処理
	defer func() {
		var errmsg string
		if err := recover(); err != nil {
			switch err := err.(type) {
			case error:
				errmsg = err.Error()
			default:
				errmsg = fmt.Sprint(err)
			}
			sends = []*discordgo.MessageSend{
				{
					Embed: &discordgo.MessageEmbed{
						Color:       0xff0000,
						Description: "``` " + strings.ReplaceAll(errmsg, "```", " `` ") + " ```",
						Title:       "不明なエラー",
					},
				},
			}
			return
		}
	}()

	error_count := 0
	send := &discordgo.MessageSend{}
	embed := &discordgo.MessageEmbed{}
	send.Embed = embed
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
				if media.LongCode || media.Error != 0 {
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
					line = "||" + strings.ReplaceAll(line, "||", "\\|\\|") + "||"
				}
				field.Value += line
			} else {
				field.Name = str
			}
		case proto.OutputMedia_FILE:
			if filepath.Ext(media.Filename) == ".svg" {
				bf := bytes.NewBuffer(media.Data)
				svg, _ := svg2jpg(bf)
				media.Data = svg.Bytes()
				media.Filename = media.Filename[:len(media.Filename)-3] + "jpg"
			}
		}
		if i+1 == len(msg.Medias) || !media.ExtendField || media.Error != 0 {
			embed.Fields = append(embed.Fields, field)
			field = &discordgo.MessageEmbedField{}
		}
		if i+1 == len(msg.Medias) || media.Type == proto.OutputMedia_FILE || media.Error != 0 {
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
			if media.Error != 0 {
				embed.Title = "Error"
				embed.Color = 0xff0000
			}
			author := &discordgo.MessageEmbedAuthor{}
			author.Name = message.Message.Author.Username
			author.IconURL = message.Author.AvatarURL("")
			embed.Author = author
			if media.Type == proto.OutputMedia_FILE {
				file := &discordgo.File{}
				send.File = file
				file.Reader = bytes.NewBuffer(media.Data)
				file.Name = media.Filename
				ext := filepath.Ext(media.Filename)
				if ext == ".png" || ext == ".jpg" || ext == ".jpeg" {
					image := &discordgo.MessageEmbedImage{}
					image.URL = fmt.Sprintf("attachment://%s", file.Name)
					embed.Image = image
				}
			}
			sends = append(sends, send)
			send = &discordgo.MessageSend{}
			embed = &discordgo.MessageEmbed{}
		}
	}
	return sends
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

func svg2jpg(svg io.Reader) (jpg *bytes.Buffer, err error) {
	icon, err := oksvg.ReadIconStream(svg)
	if err != nil {
		return nil, err
	}
	const (
		w = 800
		h = 128
	)
	icon.SetTarget(0, 0, float64(w), float64(h))
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	icon.Draw(rasterx.NewDasher(w, h, rasterx.NewScannerGV(w, h, rgba, rgba.Bounds())), 1)
	jpg = bytes.NewBuffer([]byte{})
	err = png.Encode(jpg, rgba)
	if err != nil {
		return nil, err
	}
	return
}
