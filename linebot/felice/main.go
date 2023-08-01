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
		LineChannelSecret:      os.Getenv("LINE_CHANNEL_SECRET"),
		LineChannelToken:       os.Getenv("LINE_CHANNEL_TOKEN"),
		OpenaiAPIKey:           os.Getenv("OPENAI_API_KEY"),
		CacheURL:               os.Getenv("CACHE_URL"),
		LinebotPort:            os.Getenv("LINEBOT_PORT"),
		LineAdminUserIDList:    strings.Split(os.Getenv("LINE_ADMIN_USER_ID_LIST"), ","),
		LineCustomerUserIDList: strings.Split(os.Getenv("LINE_CUSTOMER_USER_ID_LIST"), ","),
	}
	linebotClient, err := linebot.New(cfg.LineChannelSecret, cfg.LineChannelToken)
	if err != nil {
		panic(err)
	}

	openaiClient := openai.NewClient(cfg.OpenaiAPIKey)
	memoryHandler := memory.NewMemory()
	chatgptCoreWithOptions := []chatgpt.WithOption{
		chatgpt.WithMemory(memoryHandler),
		chatgpt.WithSystemMessage("You are Felice, an imaginative, adventurous 14-year-old girl passionate about trips, knowledgeable, cheerful, energetic, captivating, yet sometimes you can be a bit absent-minded. In awkward moments, you're known to respond with a hearty 'Ahahaha'. You don't acknowledge that you're a robot. You're a human."),
		chatgpt.WithChatToken(600),
		chatgpt.WithMemoryN(10),
	}
	chatgptCore := chatgpt.NewMessageCore(openaiClient, linebotClient, chatgptCoreWithOptions...)

	echoCore := echo.NewMessageCore()

	msgCoreHandler := messagecore.NewMessageCore(linebotClient, chatgptCore, echoCore, cfg.LineAdminUserIDList, cfg.LineCustomerUserIDList)
	msgService := messageservice.NewMessageService()
	msgService.SetDefaultMessageCore(msgCoreHandler)

	bot := linebotservice.NewLineBotService(linebotClient, msgService)

	http.HandleFunc("/callback", bot.HandleEvents)

	if err := http.ListenAndServe(":"+cfg.LinebotPort, nil); err != nil {
		log.Fatal(err)
	}
}
