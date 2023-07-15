package messagecore

import (
	"context"
	"time"

	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	StopMessageEvent = linebot.NewTextMessage("⏰ Timeout. Please try again.")
)

type feliceCore struct {
	linebotClient       *linebot.Client
	chatgptCore         messagecorefactory.MessageCore
	echoCore            messagecorefactory.MessageCore
	lineAdminUserIDList []string
}

func NewMessageCore(linebotClient *linebot.Client, chatgptCore, echoCore messagecorefactory.MessageCore, lineAdminUserIDList []string) messagecorefactory.MessageCore {
	core := &feliceCore{
		linebotClient:       linebotClient,
		chatgptCore:         chatgptCore,
		echoCore:            echoCore,
		lineAdminUserIDList: lineAdminUserIDList,
	}
	return core
}

func (fc *feliceCore) Process(ctx context.Context, event *linebot.Event) (linebot.SendingMessage, error) {
	userID := event.Source.UserID

	if fc.isAdmin(userID) {
		waitCh := make(chan struct{})
		var sendingMessage linebot.SendingMessage
		var err error
		go func() {
			defer close(waitCh)
			sendingMessage, err = fc.chatgptCore.Process(ctx, event)
		}()

		for {
			select {
			case <-time.After(10 * time.Second):
				return StopMessageEvent, nil
			case <-waitCh:
				return sendingMessage, err
			}
		}
	} else {
		return fc.echoCore.Process(ctx, event)
	}
}

func (fc *feliceCore) isAdmin(userID string) bool {
	for _, adminUserID := range fc.lineAdminUserIDList {
		if userID == adminUserID {
			return true
		}
	}
	return false
}
