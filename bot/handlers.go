package bot

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/H0ax0/sadcat2/config"
	"github.com/H0ax0/sadcat2/model"
	"github.com/H0ax0/sadcat2/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	SuperUser int = 800857015
)

func Sad_start(m *tb.Message) {
	Sadbot.Send(m.Sender, "started")
	fmt.Print(m.Sender.ID)
}

func sad_request_access(m *tb.Message) {
	Sadbot.Send(m.Sender, "wait for admins desicions...")
	model.NewClient(m.Sender)

	var inlineKeys [][]tb.InlineButton

	deny_inline_btn := tb.InlineButton{
		Unique: "deny_access_btn",
		Text:   "Deny",
		Data:   strconv.Itoa(m.Sender.ID),
	}

	Sadbot.Handle(&deny_inline_btn, deny_access_btn)

	allow_inline_btn := tb.InlineButton{
		Unique: "allow_access_btn",
		Text:   "Allow",
		Data:   strconv.Itoa(m.Sender.ID),
	}

	Sadbot.Handle(&allow_inline_btn, allow_access_btn)

	inlineKeys = append(inlineKeys, []tb.InlineButton{deny_inline_btn})
	inlineKeys = append(inlineKeys, []tb.InlineButton{allow_inline_btn})

	_, err := Sadbot.Send(&tb.User{ID: SuperUser}, fmt.Sprintf("User: %s \nFirts Name: %s\nSecond Name: %s \nID:%d ",
		m.Sender.Username, m.Sender.FirstName, m.Sender.LastName, m.Sender.ID), &tb.ReplyMarkup{InlineKeyboard: inlineKeys})

	if err != nil {
		fmt.Println(err)
	}

}

func deny_access_btn(c *tb.Callback) {
	Sadbot.Respond(c)
}

func allow_access_btn(c *tb.Callback) {
	uid, err := strconv.Atoi(c.Data)
	if err != nil {
		fmt.Print("eeeeorrrr")
	}
	model.Elevate(uid)
	Sadbot.Send(&tb.User{ID: uid}, "access to sadcat granted")
	Sadbot.Respond(c)
}

func Sad_OnText(m *tb.Message) {
	Sadbot.Send(m.Chat, "/help")
}

func Sad_Task(m *tb.Message) {
	Sadbot.Send(m.Chat, "asdfasdf")
}

func Sad_Log(m *tb.Message) {

	logs := utils.GetRecentLogs(config.LogBasePath, 5)
	var inlineKeys [][]tb.InlineButton
	for _, log := range logs {
		inlineBtn := tb.InlineButton{
			Unique: "log" + strings.Replace(strings.TrimSuffix(filepath.Base(log), ".log"), "-", "", -1),
			Text:   filepath.Base(log),
			Data:   filepath.Base(log),
		}
		Sadbot.Handle(&inlineBtn, bLogsInlineBtn)
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	Sadbot.Send(m.Chat, "log menu", &tb.ReplyMarkup{InlineKeyboard: inlineKeys})

}
func bLogsInlineBtn(c *tb.Callback) {
	logfile := &tb.Document{
		File:     tb.FromDisk(config.LogBasePath + c.Data),
		FileName: c.Data,
		MIME:     "text/plain",
	}
	Sadbot.Send(c.Message.Chat, logfile)
	Sadbot.Respond(c)
}
