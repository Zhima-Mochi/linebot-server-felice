package messagecore

import (
	"context"
	"log"

	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

type feliceCore struct {
	chatgptCore     messagecorefactory.MessageCore
	echoCore        messagecorefactory.MessageCore
	lineAdminUserID string
}

func NewMessageCore(chatgptCore, echoCore messagecorefactory.MessageCore, lineAdminUserID string) messagecorefactory.MessageCore {
	core := &feliceCore{
		chatgptCore:     chatgptCore,
		echoCore:        echoCore,
		lineAdminUserID: lineAdminUserID,
	}
	return core
}

func (fc *feliceCore) Process(ctx context.Context, event *linebot.Event) (linebot.SendingMessage, error) {
	userID := event.Source.UserID
	log.Printf("userID: %s", userID)
	if userID == fc.lineAdminUserID {
		return fc.chatgptCore.Process(ctx, event)
	} else {
		return fc.echoCore.Process(ctx, event)
	}
}
