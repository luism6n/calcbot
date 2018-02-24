package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/luism6n/calcbot/calc"

	"gopkg.in/telegram-bot-api.v4"
)

var (
	token *string
	debug *bool
	port  *string
)

func main() {
	bot, updates := setupBot()

	for update := range updates {
		if update.InlineQuery == nil {
			continue
		}

		query := update.InlineQuery.Query
		evaluation, err := calc.Evaluate(query)

		var result tgbotapi.InlineQueryResultArticle
		if err != nil {
			result = newInlineQueryResultArticle(fmt.Sprintf("%s", err.Error()))
		} else {
			result = newInlineQueryResultArticle(fmt.Sprintf("%s ~> %f", query, evaluation))
		}

		config := newInlineConfig(update.InlineQuery.ID, result)

		res, err := bot.AnswerInlineQuery(config)
		if err != nil {
			log.Printf("Error:\nerr: %s\nres: %+v\nquery: %s", err.Error(), res, query)
		}
	}
}

func setupBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	readCommandLineFlags()

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		die("Bot token is invalid: %s", *token)
	}

	bot.Debug = *debug

	u, err := url.Parse("https://luis-calc-bot.herokuapp.com/" + *token)
	if err != nil {
		die("Parsing webhook URL failed: %s", err.Error())
	}

	_, err = bot.SetWebhook(tgbotapi.WebhookConfig{
		URL: u,
	})
	if err != nil {
		die("Error setting webhook: %s", err.Error())
	}

	updates := bot.ListenForWebhook("/" + *token)
	log.Print("Listening at 0.0.0.0:" + *port)
	go http.ListenAndServe("0.0.0.0:"+*port, nil)

	return bot, updates
}

func readCommandLineFlags() {
	debug = flag.Bool("debug", false, "If the bot should run in debug mode")
	port = flag.String("port", "No port provided", "Port to listen for updates")
	token = flag.String("token", "No token provided", "The bot token")
	flag.Parse()

	log.Printf("Read arguments.\ndebug: %t\ntoken: %s\nport: %s", *debug, *token, *port)
}

func newInlineQueryResultArticle(text string) tgbotapi.InlineQueryResultArticle {
	return tgbotapi.InlineQueryResultArticle{
		Type:        "article",
		ID:          "only result",
		Title:       "Evaluation result",
		Description: text,
		InputMessageContent: tgbotapi.InputTextMessageContent{
			Text: text,
		},
	}
}

func newInlineConfig(queryID string, onlyResult tgbotapi.InlineQueryResultArticle) tgbotapi.InlineConfig {
	return tgbotapi.InlineConfig{
		InlineQueryID: queryID,
		Results:       castToInterfaceSlice([]tgbotapi.InlineQueryResultArticle{onlyResult}),
		CacheTime:     300,
	}
}

func castToInterfaceSlice(iqra []tgbotapi.InlineQueryResultArticle) []interface{} {
	s := make([]interface{}, len(iqra))
	for i, v := range iqra {
		s[i] = v
	}

	return s
}

func die(format string, a ...interface{}) {
	fmt.Printf(format+"\n", a...)
	os.Exit(1)
}
