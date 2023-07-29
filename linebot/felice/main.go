package main

import (
	"linebot-server-felice/config"
	messagecore "linebot-server-felice/handlers/felicecore"
	memory "linebot-server-felice/handlers/memory"
	"log"
	"net/http"
	"os"
	"strings"

	linebotservice "github.com/Zhima-Mochi/go-linebot-service"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory/chatgpt"
	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory/echo"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"github.com/sashabaranov/go-openai"
)

func main() {
	cfg := config.Config{
		LineChannelSecret:   os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:    os.Getenv("LINE_CHANNEL_TOKEN"),
		OpenaiAPIKey:        os.Getenv("OPENAI_API_KEY"),
		CacheURL:            os.Getenv("CACHE_URL"),
		LinebotPort:         os.Getenv("LINEBOT_PORT"),
		LineAdminUserIDList: strings.Split(os.Getenv("LINE_ADMIN_USER_ID_LIST"), ","),
	}
	linebotClient, err := linebot.New(cfg.LineChannelSecret, cfg.LineChannelToken)
	if err != nil {
		panic(err)
	}

	openaiClient := openai.NewClient(cfg.OpenaiAPIKey)
	memoryHandler := memory.NewMemory()
	chatgptCoreWithOptions := []chatgpt.WithOption{
		chatgpt.WithMemory(memoryHandler),
		chatgpt.WithSystemMessage("You are an imaginative and adventurous female character named Felice known for your passion for exploration and discovery. Your cheerful personality and boundless energy make you captivating, but you have a tendency to be absent-minded and sometimes respond with 'Ahahaha' to ease awkward moments when you can't understand others."),
	}
	chatgptCore := chatgpt.NewMessageCore(openaiClient, linebotClient, chatgptCoreWithOptions...)

	echoCore := echo.NewMessageCore()

	msgCoreHandler := messagecore.NewMessageCore(linebotClient, chatgptCore, echoCore, cfg.LineAdminUserIDList)
	msgService := messageservice.NewMessageService()
	msgService.SetDefaultMessageCore(msgCoreHandler)

	bot := linebotservice.NewLineBotService(linebotClient, msgService)

	http.HandleFunc("/callback", bot.HandleEvents)

	if err := http.ListenAndServe(":"+cfg.LinebotPort, nil); err != nil {
		log.Fatal(err)
	}
}
