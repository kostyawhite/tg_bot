package server

import (
	"fmt"
	"github.com/kostyawhite/telegram-bot/pkg/config"
	"github.com/kostyawhite/telegram-bot/pkg/repository"
	"github.com/zhashkevych/go-pocket-sdk"
	"net/http"
	"strconv"
)

type AuthServer struct {
	server          *http.Server
	pocketClient    *pocket.Client
	tokenRepository repository.TokenRepository
	cfg             *config.Config
	tgBotUrl        string
}

func NewAuthServer(pocketClient *pocket.Client, tokenRepository repository.TokenRepository, cfg *config.Config) *AuthServer {
	return &AuthServer{pocketClient: pocketClient, tokenRepository: tokenRepository, tgBotUrl: cfg.TgBotUrl}
}

func (as *AuthServer) Start() error {
	as.server = &http.Server{
		Addr:    ":80",
		Handler: as,
	}

	return as.server.ListenAndServe()
}

func (as *AuthServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	chatIdParam := r.URL.Query().Get("chat_id")
	if chatIdParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chatId, err := strconv.ParseInt(chatIdParam, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	requestToken, err := as.tokenRepository.Get(chatId, repository.RequestTokens)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	authResp, err := as.pocketClient.Authorize(r.Context(), requestToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = as.tokenRepository.Save(chatId, authResp.AccessToken, repository.AccessTokens)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Printf("chat_id: %d\nrequest_token: %s\naccess_token: %s\n", chatId, requestToken, authResp.AccessToken)

	w.Header().Add("Location", as.tgBotUrl)
	w.WriteHeader(http.StatusMovedPermanently)
}
