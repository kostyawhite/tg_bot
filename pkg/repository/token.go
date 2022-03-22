package repository

type Bucket string

const (
	AccessTokens  Bucket = "access_tokens"
	RequestTokens Bucket = "request_tokens"
)

type TokenRepository interface {
	Save(chatId int64, token string, bucket Bucket) error
	Get(chatId int64, bucket Bucket) (string, error)
}
