package messagecore

import (
	"context"
	"linebot-server-felice/handlers/cache"
	"log"
	"time"

	"github.com/Zhima-Mochi/go-linebot-service/messageservice/messagecorefactory"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

var (
	StopMessageEvent = linebot.NewTextMessage("⏰ Timeout. Please try again.")
)

type feliceCore struct {
	linebotClient          *linebot.Client
	chatgptCore            messagecorefactory.MessageCore
	echoCore               messagecorefactory.MessageCore
	lineAdminUserIDList    []string
	lineCustomerUserIDList []string
	cacheHandler           *cache.CacheHandler
}

func NewMessageCore(linebotClient *linebot.Client, chatgptCore, echoCore messagecorefactory.MessageCore, lineAdminUserIDList, lineCustomerUserIDList []string) messagecorefactory.MessageCore {
	core := &feliceCore{
		linebotClient:          linebotClient,
		chatgptCore:            chatgptCore,
		echoCore:               echoCore,
		lineAdminUserIDList:    lineAdminUserIDList,
		lineCustomerUserIDList: lineCustomerUserIDList,
		cacheHandler:           cache.NewCacheHandler(),
	}
	return core
}

func (fc *feliceCore) Process(ctx context.Context, event *linebot.Event) (linebot.SendingMessage, error) {
	userID := event.Source.UserID

	waitCh := make(chan struct{})
	var sendingMessage linebot.SendingMessage
	var err error
	if fc.isAdmin(userID) {
		go func() {
			defer close(waitCh)
			sendingMessage, err = fc.chatgptCore.Process(ctx, event)
		}()

		for {
			select {
			case <-time.After(30 * time.Second):
				return StopMessageEvent, nil
			case <-waitCh:
				return sendingMessage, err
			}
		}
	} else if fc.isCustomer(userID) {
		if ok, err := fc.cacheHandler.SetNX(ctx, "customer:"+userID, "1", 15*time.Second); err != nil {
			log.Fatal(err)
			return nil, err
		} else if !ok {
			return fc.echoCore.Process(ctx, event)
		} else {
			go func() {
				defer close(waitCh)
				sendingMessage, err = fc.chatgptCore.Process(ctx, event)
			}()

			for {
				select {
				case <-time.After(30 * time.Second):
					return StopMessageEvent, nil
				case <-waitCh:
					return sendingMessage, err
				}
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

func (fc *feliceCore) isCustomer(userID string) bool {
	for _, customerUserID := range fc.lineCustomerUserIDList {
		if userID == customerUserID {
			return true
		}
	}
	return false
}
