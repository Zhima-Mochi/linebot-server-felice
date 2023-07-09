package main

import (
	"linebot-server-felice/config"
	"linebot-server-felice/handlers/messagecore"
	"log"
	"net/http"
	"os"

	linebotservice "github.com/Zhima-Mochi/go-linebot-service"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory/chatgpt"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory/echo"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sashabaranov/go-openai"
)

var (
	cfg = config.Config{}
)

func init() {
	cfg = config.Config{
		LineChannelSecret: os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:  os.Getenv("LINE_CHANNEL_TOKEN"),
		OpenaiAPIKey:      os.Getenv("OPENAI_API_KEY"),
		CacheURL:          os.Getenv("CACHE_URL"),
		LinebotPort:       os.Getenv("LINEBOT_PORT"),
		LineAdminUserID:   os.Getenv("LINE_ADMIN_USER_ID"),
	}

}

func main() {
	client, err := linebot.New(cfg.LineChannelSecret, cfg.LineChannelToken)
	if err != nil {
		panic(err)
	}

	openaiClient := openai.NewClient(cfg.OpenaiAPIKey)
	chatgptCore := chatgpt.NewMessageCore(openaiClient)

	echoCore := echo.NewMessageCore()

	msgCoreHandler := messagecore.NewMessageCore(chatgptCore, echoCore)
	msgService := messageservice.NewMessageService()
	msgService.SetDefaultMessageCore(msgCoreHandler)

	bot := linebotservice.NewLineBotService(client, msgService)

	http.HandleFunc("/callback", bot.Do)

	if err := http.ListenAndServe(":"+cfg.LinebotPort, nil); err != nil {
		log.Fatal(err)
	}
}
