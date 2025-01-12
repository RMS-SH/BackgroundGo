package validators

import (
	"fmt"
	"regexp"
	"strings"
)

// MessageValidator é responsável por validar e filtrar mensagens.
type MessageValidator struct {
	PalavrasProibidas []string
	FilterRegex       *regexp.Regexp
}

// NewMessageValidator cria uma nova instância de MessageValidator.
func NewMessageValidator() *MessageValidator {
	return &MessageValidator{
		PalavrasProibidas: []string{},
		FilterRegex:       nil,
	}
}

// SetFilterByRegex compila a expressão regular fornecida e a define como filtro.
func (v *MessageValidator) SetFilterByRegex(pattern string) error {
	compiledRegex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("erro ao compilar regex: %v", err)
	}
	v.FilterRegex = compiledRegex
	return nil
}

// SetPalavrasProibidas define a lista de palavras proibidas.
func (v *MessageValidator) SetPalavrasProibidas(palavras []string) {
	v.PalavrasProibidas = palavras
}

// DeveEnviarMensagem verifica se a mensagem pode ser enviada.
// Retorna false se a mensagem contiver qualquer palavra proibida.
func (v *MessageValidator) DeveEnviarMensagem(mensagem string) bool {
	for _, palavra := range v.PalavrasProibidas {
		if strings.Contains(mensagem, palavra) {
			return false
		}
	}
	return true
}

// ApplyFilterByRegex aplica a expressão regular definida para filtrar a mensagem.
// Remove todas as ocorrências que correspondem ao padrão regex.
func (v *MessageValidator) ApplyFilterByRegex(mensagem string) string {
	if v.FilterRegex == nil {
		return mensagem
	}
	// Substitui todas as ocorrências que correspondem ao regex por uma string vazia
	cleanedText := v.FilterRegex.ReplaceAllString(mensagem, "")
	// Remove espaços em branco no início e no fim
	cleanedText = strings.TrimSpace(cleanedText)
	return cleanedText
}
