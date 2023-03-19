package bot

import (
	"log"
	"net/http"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/slonoboy/go-bot/internal/database"
)

const (
	currencies       = "currencies"
	cryptocurrencies = "cryptocurrencies"
)

type Bot struct {
	telegram *tgbotapi.BotAPI
	api      *resty.Client
	dh       *database.DatabaseHandler
}

func NewBot(token string, api *resty.Client, dh *database.DatabaseHandler, logger *logrus.Logger) (*Bot, error) {
	tg, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	api.SetLogger(logger)

	return &Bot{
		telegram: tg,
		api:      api,
		dh:       dh,
	}, nil
}

func (b *Bot) StartWebHook(url string) error {
	webhookURL := url + "/" + b.telegram.Token
	_, err := tgbotapi.NewWebhook(webhookURL)
	if err != nil {
		return err
	}

	http.HandleFunc("/"+b.telegram.Token, func(w http.ResponseWriter, r *http.Request) {
		update, err := b.telegram.HandleUpdate(r)
		if err != nil {
			log.Printf("Error handling update: %s", err.Error())
			return
		}
		b.handleUpdate(update)
	})

	return nil
}

func (b *Bot) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.telegram.GetUpdatesChan(u)

	for update := range updates {
		if err := b.handleUpdate(&update); err != nil {
			log.Println(err)
		}
	}
}

func (b *Bot) handleUpdate(update *tgbotapi.Update) error {
	switch {
	case update.Message != nil && update.Message.IsCommand():
		return b.handleCommand(update)
	case update.CallbackQuery != nil:
		return b.handleCallBackQuery(update.CallbackQuery)
	default:
		return nil
	}
}

func (b *Bot) handleCommand(update *tgbotapi.Update) error {
	switch update.Message.Command() {
	case "start":
		return b.handleStart(update)
	case "help":
		return b.handleHelp(update)
	case "convert":
		return b.handleConvert(update)
	case "my_cryptos":
		return b.handleMyMostFrequentCryptosOrCurrencies(update, cryptocurrencies)
	case "my_currencies":
		return b.handleMyMostFrequentCryptosOrCurrencies(update, currencies)
	case "my_first_request_date":
		return b.handleMyFirstRequestDate(update)
	case "my_requests_num":
		return b.handleMyRequestsNum(update)
	default:
		return nil
	}
}

func (b *Bot) handleCallBackQuery(q *tgbotapi.CallbackQuery) error {
	return nil
}
