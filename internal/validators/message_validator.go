package validators

import "strings"

type MessageValidator struct {
	PalavrasProibidas []string
}

func NewMessageValidator() *MessageValidator {
	return &MessageValidator{
		PalavrasProibidas: []string{},
	}
}

func (v *MessageValidator) SetPalavrasProibidas(palavras []string) {
	v.PalavrasProibidas = palavras
}

func (v *MessageValidator) DeveEnviarMensagem(mensagem string) bool {
	for _, palavra := range v.PalavrasProibidas {
		if strings.Contains(mensagem, palavra) {
			return false
		}
	}
	return true
}
