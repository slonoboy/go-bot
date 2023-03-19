package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/slonoboy/go-bot/config"
	botHandler "github.com/slonoboy/go-bot/internal/bot"
	"github.com/slonoboy/go-bot/internal/database"
	"github.com/slonoboy/go-bot/internal/server"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	cfg, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		cfg.Database.Host,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}

	dh, err := database.NewDatabaseHandler(db)
	if err != nil {
		log.Panic(err)
	}

	logger := setUpLogger(*cfg)
	bot, err := botHandler.NewBot(cfg.Bot.Token, resty.New(), dh, logger)
	if err != nil {
		log.Panic(err)
	}

	srv := server.NewServer(cfg, nil)

	bot.StartWebHook(cfg.HTTP.Domain)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatalf("Произошла ошибка во время запуска сервера: %s\n", err.Error())
		}
	}()

	log.Print("Сервер запущен")

	<-make(chan struct{})
}

func setUpLogger(cfg config.Config) *logrus.Logger {
	logger := logrus.New()

	logFilePath := "/logs/logs.log"
	logRotation := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100, // 100 MB
		MaxBackups: 30,
		LocalTime:  true,
		Compress:   true,
	}

	mw := io.MultiWriter(os.Stdout, logRotation)
	logger.SetOutput(mw)

	logger.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: time.RFC3339,
	})
	return logger
}
