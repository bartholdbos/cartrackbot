package main

import(
	"github.com/bartholdbos/golegram"
	"log"
	"fmt"
	"net/http"
)

func handleErr(err error) {
	log.Fatal(err)
}

func handlePing(out http.ResponseWriter, in *http.Request) {
	out.Write([]byte("pong"))
}

func handleUpdate(update golegram.Update) {
	fmt.Println(update.Message.Text)
}

func main() {
	bot, err := golegram.NewBot("147357073:AAHTaz9TsAY0SbVXAEyfuYeRJGNN2m_hwwk")

	if err != nil {
		handleErr(err)
	}

	bot.AddToWebhook(handleUpdate, handlePing)
}