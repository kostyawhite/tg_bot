package telegram

import (
	"context"
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kostyawhite/telegram-bot/pkg/config"
	"github.com/kostyawhite/telegram-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/url"
)

type Bot struct {
	bot             *tgbotapi.BotAPI
	pocketClient    *pocket.Client
	authServerLink  string
	tokenRepository repository.TokenRepository
	messages        config.Messages
}

func NewBot(bot *tgbotapi.BotAPI, client *pocket.Client, authServerLink string, tokenRepository repository.TokenRepository, messages config.Messages) *Bot {
	return &Bot{bot: bot, pocketClient: client, authServerLink: authServerLink, tokenRepository: tokenRepository, messages: messages}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	b.handleUpdates(updates)
	return nil
}

const (
	commandStart = "start"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	command := message.Command()

	switch command {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return errUnknownCommand
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	_, err := b.getAccessToken(message.Chat.ID)
	if err != nil {
		return b.initAuthProcess(message)
	}
	msg := tgbotapi.NewMessage(message.Chat.ID, b.messages.ReplyAlreadyAuthorized)
	_, err = b.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) handleReplyMessage(message *tgbotapi.Message) {
	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	b.bot.Send(msg)
}

func (b *Bot) handleNewMessage(message *tgbotapi.Message) error {
	chatId := message.Chat.ID
	msg := tgbotapi.NewMessage(chatId, b.messages.LinkSavedSuccessfully)
	_, err := url.ParseRequestURI(message.Text)
	if err != nil {
		return errInvalidUrl
	}

	accessToken, err := b.getAccessToken(chatId)
	if err != nil {
		return errUnauthorized
	}

	if err = b.pocketClient.Add(context.Background(), pocket.AddInput{
		AccessToken: accessToken,
		URL:         message.Text,
	}); err != nil {
		return errUnableToSave
	}

	_, err = b.bot.Send(msg)
	return err
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message != nil { // If we got a message

			if update.Message.IsCommand() {
				if err := b.handleCommand(update.Message); err != nil {
					b.handleError(update.Message.Chat.ID, err)
				}
				continue
			}

			if err := b.handleNewMessage(update.Message); err != nil {
				b.handleError(update.Message.Chat.ID, err)
			}
		}
	}
}

var (
	errInvalidUrl     = errors.New("invalid url")
	errUnableToSave   = errors.New("unable to save link")
	errUnauthorized   = errors.New("unauthorized")
	errUnknownCommand = errors.New("unknown command")
)

func (b Bot) handleError(chatId int64, err error) {
	msg := tgbotapi.NewMessage(chatId, "")
	switch err {
	case errInvalidUrl:
		msg.Text = b.messages.InvalidLink
		b.bot.Send(msg)
	case errUnableToSave:
		msg.Text = b.messages.FailedToSave
		b.bot.Send(msg)
	case errUnauthorized:
		msg.Text = b.messages.Unauthorized
		b.bot.Send(msg)
	case errUnknownCommand:
		msg.Text = b.messages.UnknownCommand
		b.bot.Send(msg)
	default:
		msg.Text = b.messages.UnknownError
		b.bot.Send(msg)
	}
}
