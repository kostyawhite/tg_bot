package main

import (
	"github.com/boltdb/bolt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kostyawhite/telegram-bot/pkg/config"
	"github.com/kostyawhite/telegram-bot/pkg/repository"
	"github.com/kostyawhite/telegram-bot/pkg/repository/boltdb"
	"github.com/kostyawhite/telegram-bot/pkg/server"
	"github.com/kostyawhite/telegram-bot/pkg/telegram"
	"github.com/zhashkevych/go-pocket-sdk"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Panic(err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	pocketClient, err := pocket.NewClient(cfg.ConsumerKey)
	if err != nil {
		log.Panic(err)
	}

	db, err := initDB(cfg)
	if err != nil {
		log.Panic(err)
	}

	tokenRepository := boltdb.NewTokenRepository(db)

	tgBot := telegram.NewBot(bot, pocketClient, cfg.AuthServerUrl, tokenRepository, cfg.Messages)

	authServer := server.NewAuthServer(pocketClient, tokenRepository, cfg)

	go func() {
		if err := authServer.Start(); err != nil {
			log.Panic(err)
		}
	}()

	if err := tgBot.Start(); err != nil {
		log.Panic(err)
	}
}

func initDB(cfg *config.Config) (*bolt.DB, error) {
	db, err := bolt.Open(cfg.DbPath, 0644, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(repository.AccessTokens))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte(repository.RequestTokens))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
