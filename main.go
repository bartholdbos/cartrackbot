package main

import (
	"database/sql"
	"fmt"
	"github.com/bartholdbos/golegram"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var db *sql.DB
var bot *golegram.Bot
var config Configuration

func parse(input string) (string, []string) {
	s := strings.Split(input, " ")
	command := strings.Split(s[0], "@")[0]
	var args []string
	for i := 1; i < len(s); i++ {
		if s[i] != "" {
			args = append(args, s[i])
		}
	}
	return command, args
}

func getCar(userid int32) (car int, err error) {
	row := db.QueryRow("SELECT `car` FROM `users` WHERE `id` = ?", userid)
	err = row.Scan(&car)

	return
}

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
		command, _ := parse(update.Message.Text)

		switch command {
		case "/start":
			_, err := db.Exec("INSERT INTO `users` (`id`) VALUES (?);", update.Message.From.Id)
			if err != nil{
				log.Printf("%s\n", err) // TODO: Return error to user
				break
			}
			
			bot.SendMessage(strconv.Itoa(int(update.Message.Chat.Id)), "The Cartrackbot", false, "")
		case "/location":
			car, err := getCar(update.Message.From.Id)
			if err != nil {
				log.Printf("%s\n", err) // TODO: Return error to user
				break
			}

			if car == 0 {
				bot.SendMessage(strconv.Itoa(int(update.Message.Chat.Id)), "Je hebt nog geen auto toegevoegd", false, "")
				break
			}

			var longitude float64
			var latitude float64

			row := db.QueryRow("SELECT `longitude`, `latitude` FROM `cars` WHERE `id` = ?;", car)
			err1 := row.Scan(&longitude, &latitude)
			if err1 != nil {
				log.Printf("%s\n", err) // TODO: Return error to user
				break
			}

			bot.SendLocation(strconv.Itoa(int(update.Message.Chat.Id)), latitude, longitude, false, 0)
		}
	} else {
		car, err := getCar(update.Message.From.Id)
		if err != nil {
			log.Printf("%s\n", err) // TODO: Return error to user
			break
		}

		if car == 0 {
			bot.SendMessage(strconv.Itoa(int(update.Message.Chat.Id)), "Je hebt nog geen auto toegevoegd", false, "")
			break
		}

		_, err1 := db.Exec("UPDATE `cars` SET `longitude`=?,`latitude`=? WHERE `id` = ?;", update.Message.Location.Longitude, update.Message.Location.Latitude, car)
		if err1 != nil {
			log.Printf("%s\n", err) // TODO: Return error to user
			break
		}

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
