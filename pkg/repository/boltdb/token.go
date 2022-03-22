package boltdb

import (
	"errors"
	"github.com/boltdb/bolt"
	"github.com/kostyawhite/telegram-bot/pkg/repository"
	"strconv"
)

type TokenRepository struct {
	db *bolt.DB
}

func NewTokenRepository(db *bolt.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (t TokenRepository) Save(chatId int64, token string, bucket repository.Bucket) error {
	err := t.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(intToBytes(chatId), []byte(token))
		return err
	})
	return err
}

func (t TokenRepository) Get(chatId int64, bucket repository.Bucket) (string, error) {
	var token string

	err := t.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		data := b.Get(intToBytes(chatId))
		token = string(data)
		return nil
	})

	if err != nil {
		return "", err
	}

	if token == "" {
		return "", errors.New("token not found")
	}

	return token, nil
}

func intToBytes(v int64) []byte {
	return []byte(strconv.FormatInt(v, 10))
}
