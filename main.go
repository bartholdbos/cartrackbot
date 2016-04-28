package main

import(
	"github.com/bartholdbos/golegram"
	"log"
	"fmt"
	"net/http"
)

var bot golegram.Bot
var config Configuration

func handleErr(err error) {
	log.Fatal(err)
}

func handlePing(out http.ResponseWriter, in *http.Request) {
	out.Write([]byte("pong"))
}

func handleUpdate(update golegram.Update) {
	fmt.Println(update.Message.Text)
}

func setupConfig() {
	var err error
	config, err = load_config()

	if err != nil {
		handleErr(err)
	}

	fmt.Println("Config loaded")
}

func setupBot() {
	var err error
	bot, err := golegram.NewBot(config.TelegramBotToken)

	if err != nil {
		handleErr(err)
	}

	fmt.Println("Bot created: " + bot.User.Username)
}

func setupWebhook() {
	bot.AddToWebhook(handleUpdate, handlePing)

	fmt.Println("Starting Webhook")
	err := golegram.StartWebhook(config.Webhook.Port, config.Webhook.Public, config.Webhook.Private)

	if err != nil {
		handleErr(err)
	}
}

func main() {
	setupConfig()
	setupBot()
	setupWebhook()
}