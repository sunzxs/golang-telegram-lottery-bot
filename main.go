package main

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/telegram-bot-api.v4"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

const (
	TIMEOUT       = 60
	DATABASE      = "sqlite3"
	DATABASE_NAME = "db.sqlite3"
)

var (
	bot *tgbotapi.BotAPI
	gdb *gorm.DB
)

// You must create bot_token.go file, which include TOKEN variable in global package scope
func init() {
	var err error

	gdb, err = gorm.Open(DATABASE, DATABASE_NAME)
	if err != nil {
		panic(err)
	}
	gdb.LogMode(true)

	gdb.AutoMigrate(
		&User{},
		&Info{},
	)

	bot, err = tgbotapi.NewBotAPI(TOKEN)

	if err != nil {
		err.Error()
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

}

func main() {
	defer gdb.Close()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = TIMEOUT

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		err.Error()
	}

	for update := range updates {
		msg := update.Message
		if msg == nil {
			continue
		} else if msg.IsCommand() {
			command := msg.Command()

			switch command {
				// Регистрация
			case "start":
				reg(msg, update)
				// Старт лотореи
			case "begin":
				if checkAdminAccess(msg) {
					start(msg, update)
				}
				// Стоп
			case "finish":
				if checkAdminAccess(msg) {
					stop(msg, update)
				}
				// Список участников
			case "list":
				if checkAdminAccess(msg) {
					list(msg)
				}
				// Разыграть
			case "startLottery":
				if checkAdminAccess(msg) {
					startLottery(msg)
				}
				// Сообщение победителям
			case "winners":
				if checkAdminAccess(msg) {
					messageToWinners(msg)
				}
			case "regstop":
				if checkAdminAccess(msg) {
					regstop(msg)
				}
			case "regstart":
				if checkAdminAccess(msg) {
					regstart(msg)
				}
			}
		}
	}
}
