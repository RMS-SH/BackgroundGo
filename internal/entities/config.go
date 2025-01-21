package entities

type Config struct {
	ApiKeyOpenIA  string
	ApiKeyEntrega string
	PromptArquivo string
	PromptImagem  string
	UUIDUser      string
	WorkSpaceID   string
	ApiKeyBot     string
	BaseUrlUchat  string
	URLMotorIA    string
	Telefone      string
	Nome          string
	ExtraVars     map[string]string
}

func NewConfig(
	MessageItem MessageItem,
	ApiOpenIA string,
	BaseUrlUchat string,
	ExtraVars map[string]string,
) Config {
	return Config{
		ApiKeyOpenIA:  ApiOpenIA,
		PromptArquivo: MessageItem.PromptArquivo,
		PromptImagem:  MessageItem.PromptImagem,
		UUIDUser:      MessageItem.UserNS,
		WorkSpaceID:   MessageItem.NomeWorkspace,
		ApiKeyBot:     MessageItem.ApiKeyBot,
		BaseUrlUchat:  BaseUrlUchat,
		URLMotorIA:    MessageItem.MotorIA,
		Telefone:      MessageItem.Telefone,
		Nome:          MessageItem.Nome,
		ExtraVars:     ExtraVars,
	}
}
