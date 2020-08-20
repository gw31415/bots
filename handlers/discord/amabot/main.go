package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gw31415/bots/lib/discord"
	"github.com/gw31415/bots/lib"
	"github.com/joho/godotenv"
)

func main() {
	// .env
	if len(os.Args) == 1 {
		if err := godotenv.Load(); err != nil {
			fmt.Println("info: " + err.Error())
		}
	} else {
		if err := godotenv.Load(os.Args[1:]...); err != nil {
			panic(err)
		}
	}

	// Instanceを作成
	instance, err := discord.NewInstance(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}
	// botsのコマンドハンドラを作成
	wd, _ := os.Getwd()
	bots, err := lib.NewBotsHandler(wd + "/bin")
	if err != nil {
		panic(err)
	}
	// UnaryServiceの設定を構成
	preprocessor := discord.Preprocessor_prefix(os.Getenv("PREFIX"))
	style := discord.MessageEmbedStyle_default
	config := discord.NewUnaryServiceConfig(preprocessor, style)
	// UnaryServiceを設定
	instance.SetUnaryService(bots, config)

	//Instanceを起動
	instance.Open()

	// SIGINT待ち
	fmt.Println("Press Ctrl-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Instanceを停止
	instance.Close()
}
