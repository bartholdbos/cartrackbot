package main

import (
	"database/sql"
	"fmt"
	"github.com/bartholdbos/golegram"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

var db *sql.DB
var bot *golegram.Bot
var config Configuration

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handlePing(out http.ResponseWriter, in *http.Request) {
	out.Write([]byte("pong"))
}

func handleUpdate(update golegram.Update) {
	if update.Message.Location.Longitude == 0 && update.Message.Location.Latitude == 0 {
		bot.SendMessage(strconv.Itoa(int(update.Message.Chat.Id)), "test", false, "")
	} else {
		row, err := db.Query("SELECT `car` FROM `users` WHERE `id` = ?;", update.Message.Chat.Id)
		handleErr(err)
		row.Next()

		var car int
		err1 := row.Scan(&car)
		handleErr(err1)

		_, err2 := db.Query("UPDATE `cars` SET `longitude`=?,`latitude`=? WHERE `id` = ?;", update.Message.Location.Longitude, update.Message.Location.Latitude, car)
		handleErr(err2)

		bot.SendMessage(strconv.Itoa(int(update.Message.Chat.Id)), "Locatie bijgewerkt", false, "")
	}
}

func setupConfig() {
	var err error
	config, err = load_config()

	handleErr(err)

	fmt.Println("Config loaded")
}

func setupBot() {
	var err error
	bot, err = golegram.NewBot(config.TelegramBotToken)

	handleErr(err)

	fmt.Println("Bot created: " + bot.User.Username)
}

func setupMySQL() {
	var err error
	db, err = sql.Open("mysql", config.generateDSN())

	handleErr(err)
	handleErr(db.Ping())
}

func setupWebhook() {
	bot.AddToWebhook(handleUpdate, handlePing)

	fmt.Println("Starting Webhook")
	err := golegram.StartWebhook(config.Webhook.Port, config.Webhook.Public, config.Webhook.Private)

	handleErr(err)
}

func main() {
	setupConfig()
	setupBot()
	setupMySQL()
	setupWebhook()
}
