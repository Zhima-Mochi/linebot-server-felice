package messagecore

import (
	"linebot-server-felice/config"
	"os"

	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	cfg config.Config
)

func init() {
	cfg = config.Config{
		LineAdminUserID: os.Getenv("LINE_ADMIN_USER_ID"),
	}
}

type feliceCore struct {
	chatgptCore messagecorefactory.MessageCore
	echoCore    messagecorefactory.MessageCore
}

func NewMessageCore(chatgptCore, echoCore messagecorefactory.MessageCore) messagecorefactory.MessageCore {
	core := &feliceCore{
		chatgptCore: chatgptCore,
		echoCore:    echoCore,
	}
	return core
}

func (fc *feliceCore) Process(event *linebot.Event) (linebot.SendingMessage, error) {
	userID := event.Source.UserID
	if userID == cfg.LineAdminUserID {
		return fc.chatgptCore.Process(event)
	} else {
		return fc.echoCore.Process(event)
	}
}
