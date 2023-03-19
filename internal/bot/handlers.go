package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/slonoboy/go-bot/internal/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleStart(update *tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bot Started")
	if _, err := b.telegram.Send(msg); err != nil {
		log.Panic(err)
	}
	tgUser := update.SentFrom()
	user := database.User{
		TGID:      tgUser.ID,
		Username:  tgUser.UserName,
		FirstName: tgUser.FirstName,
		LastName:  tgUser.LastName,
	}

	err := b.dh.CreateUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleHelp(update *tgbotapi.Update) error {
	text := "Список команд:\n"
	text += "/convert <крипта:обязательно> <валюта:обязательно> конвертирует критовалюту в обычную валюту\n"
	text += "/my_cryptos <n:необязательно> возвращает ваши самые частые криптовалюты\n"
	text += "/my_currencies <n:необязательно> возвращает ваши самые частые валюты\n"
	text += "/my_first_request_date возвращает дату первого запроса и сам запрос\n"
	text += "/my_requests_num возвращает количество сделанных запросов"
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	if _, err := b.telegram.Send(msg); err != nil {
		log.Panic(err)
	}
	return nil
}

func (b *Bot) handleConvert(update *tgbotapi.Update) error {

	arguments := strings.Split(update.Message.CommandArguments(), " ")

	// Проверка наличия обоих аргументов
	if len(arguments) != 2 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Аргументы введены неправильно, повторите снова")
		if _, err := b.telegram.Send(msg); err != nil {
			return err
		}
		return nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Перевожу валюту...")
	if _, err := b.telegram.Send(msg); err != nil {
		return err
	}

	// Инициализация параметров для запроса
	params := map[string]string{
		"ids":           arguments[0],
		"vs_currencies": arguments[1],
	}

	// Вызов запроса с параметрами
	resp, err := b.api.R().SetQueryParams(params).Get("https://api.coingecko.com/api/v3/simple/price")
	if err != nil {
		return err
	}

	// Проверка успешности запроса
	if resp.Status() != "200 OK" {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Не получилось получить данные")
		if _, err := b.telegram.Send(msg); err != nil {
			return err
		}
		return nil
	}

	// Сохранение запроса в бд
	user, err := b.dh.FindUserByTGID(update.SentFrom().ID)
	if err != nil {
		return err
	}
	conversion := database.Conversion{
		UserID:   user.ID,
		ChatID:   update.Message.Chat.ID,
		Crypto:   arguments[0],
		Currency: arguments[1],
	}

	err = b.dh.AddConversion(conversion)
	if err != nil {
		return err
	}

	// Генерация сообщения на основе json ответа
	var response map[string]map[string]float32
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, makeResponseReadable(response))
	if _, err := b.telegram.Send(msg); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleMyMostFrequentCryptosOrCurrencies(update *tgbotapi.Update, cType string) error {
	argument := update.Message.CommandArguments()
	n, err := strconv.Atoi(argument)
	if err != nil {
		n = 3
	}

	var counts []database.ConversionCount
	if cType == cryptocurrencies {
		counts, err = b.dh.UserMostFrequentCryptos(n)
	} else if cType == currencies {
		counts, err = b.dh.UserMostFrequentCurrencies(n)
	} else {
		return errors.New("currency type not found")
	}

	if err != nil {
		return err
	}

	var result string
	for i := range counts {
		result += fmt.Sprintf("%v) %s%s: %v\n", i+1, counts[i].Crypto, counts[i].Currency, counts[i].Count)
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	if _, err := b.telegram.Send(msg); err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleMyFirstRequestDate(update *tgbotapi.Update) error {
	conversion, err := b.dh.UserFirstConversion(update.SentFrom().ID)
	if err != nil {
		return err
	}

	dateAndTime := conversion.CreatedAt.Format("02:01:2006 15:04:05")
	result := fmt.Sprintf("Дата и время первого запроса: %s", dateAndTime)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	if _, err := b.telegram.Send(msg); err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleMyRequestsNum(update *tgbotapi.Update) error {
	count, err := b.dh.UserConversionsCount(update.SentFrom().ID)
	if err != nil {
		return err
	}

	result := fmt.Sprintf("Общее количество запросов: %v", count)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, result)
	if _, err := b.telegram.Send(msg); err != nil {
		return err
	}
	return nil
}

func makeResponseReadable(response map[string]map[string]float32) string {
	var resultString = ""

	for crypto, vs_currencies := range response {
		for currency, price := range vs_currencies {
			resultString += fmt.Sprintf("%s -> %s = %v\n", crypto, currency, price)
		}
		resultString += "\n"
	}

	return resultString
}
