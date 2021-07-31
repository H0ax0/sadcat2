package bot

import (
	"fmt"
	"strconv"
	"time"

	"github.com/H0ax0/sadcat2/config"
	"github.com/H0ax0/sadcat2/model"
	"go.uber.org/zap"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	Sadbot *tb.Bot
)

func BotStart() {
	var err error

	Poller := &tb.LongPoller{Timeout: 15 * time.Second}
	spamPoller := tb.NewMiddlewarePoller(Poller, func(upd *tb.Update) bool {

		if upd.Message == nil {
			return true
		}

		if !upd.Message.Private() {
			return false
		}

		if !model.IsAuth(upd.Message.Sender.ID) && upd.Message.Text != "/request_access" {
			_, err := Sadbot.Reply(upd.Message, "you are not autorized to send commands,use /request_access")
			if err != nil {
				fmt.Println("fail")
			}
			return false
		}

		return true
	})

	botSetting := tb.Settings{
		Token:  config.Token,
		Poller: spamPoller,
	}

	//create bot
	Sadbot, err = tb.NewBot(botSetting)
	if err != nil {
		zap.S().Errorw("failed to create bot", "error", err)
		return
	}
	fmt.Println("Bot: " + strconv.Itoa(Sadbot.Me.ID) + " " + Sadbot.Me.Username)

	MakeHandle()
	fmt.Println("Bot Start")
	fmt.Println("------------")
	Sadbot.Start()
}

func MakeHandle() {

	Sadbot.Handle("/request_access", sad_request_access)
	Sadbot.Handle(tb.OnText, Sad_OnText)
	Sadbot.Handle("/task", Sad_Task)
	Sadbot.Handle("/log", Sad_Log)
}
