package main

import (
	tgapi "github.com/StepanShevelev/tg-bot-bot/pkg/api"
	cfg "github.com/StepanShevelev/tg-bot-bot/pkg/config"
	mydb "github.com/StepanShevelev/tg-bot-bot/pkg/db"
	tlg "github.com/StepanShevelev/tg-bot-bot/pkg/telegram"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
)

func main() {

	c := make(chan os.Signal, 1)

	config := cfg.New()
	if err := config.Load("./configs", "config", "yml"); err != nil {
		log.Fatal(err)
	}

	mydb.ConnectToDb()

	bot, update := tlg.BotInit(config)
	exitCh := make(chan struct{}, 1)
	go tlg.BotServe(bot, update, exitCh)
	tgapi.InitBackendApi()
	http.ListenAndServe(":"+config.Port, nil)
	signal.Notify(c, os.Interrupt)
	<-c
	exitCh <- struct{}{}
}
