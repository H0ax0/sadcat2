package main

import (
	"github.com/H0ax0/sadcat2/bot"
	"github.com/H0ax0/sadcat2/config"
	"github.com/H0ax0/sadcat2/logger"
	"github.com/H0ax0/sadcat2/model"
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	model.InitDB()
	bot.BotStart()

}
