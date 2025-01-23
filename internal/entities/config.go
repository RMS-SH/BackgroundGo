package entities

import "time"

type Config struct {
	ApiKeyOpenIA   string
	ApiKeyEntrega  string
	PromptArquivo  string
	PromptImagem   string
	UUIDUser       string
	WorkSpaceID    string
	ApiKeyBot      string
	BaseUrlUchat   string
	URLMotorIA     string
	Telefone       string
	Nome           string
	BearerTokenRMS string
	ExtraVars      map[string]string
	Timeout        time.Duration
}

func NewConfig(
	MessageItem MessageItem,
	ApiOpenIA string,
	BaseUrlUchat string,
	ExtraVars map[string]string,
	TokenRMS string,
) Config {
	return Config{
		ApiKeyOpenIA:   ApiOpenIA,
		PromptArquivo:  MessageItem.PromptArquivo,
		PromptImagem:   MessageItem.PromptImagem,
		UUIDUser:       MessageItem.UserNS,
		WorkSpaceID:    MessageItem.IDWorkSpace,
		ApiKeyBot:      MessageItem.ApiKeyBot,
		BaseUrlUchat:   BaseUrlUchat,
		URLMotorIA:     MessageItem.MotorIA,
		Telefone:       MessageItem.Telefone,
		Nome:           MessageItem.Nome,
		ExtraVars:      ExtraVars,
		BearerTokenRMS: TokenRMS,
	}
}
