package interfaces

import (
	ad "github.com/RMS-SH/BackgroundGo/adapters"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
)

type Process interface {
	ProcessaTexto(texto entities.MessageItem) (string, error)
	ProcessaImagem(texto entities.MessageItem) (string, error)
	ProcessaAudio(texto entities.MessageItem) (string, error)
	ProcessaFile(texto entities.MessageItem) (string, error)
}

type Entrega interface {
	EnviaMensagem(texto string) error
}

type DB interface {
	ArmazenaReferenciasDeArquivos(texto, url string) error
	ConsultaURLReferencia(url string) (string, error)
	ConsultaDadosEmpresa(workSpaceID string) (*entities_db.Empresa, error)
	ContagemDeRespostas(workSpaceID string) error
	ContagemDeMinutosAudio(minutos float64, workSpaceID string) error
	ContagemDeImagensProcessadas(workSpaceID string) error
	ContagemDeArquivosProcessados(workSpaceID string) error
}

type Internal interface {
	ConsultaReferenciaConversa(referenceId string) (*ad.ReferenceReturn, error)
}

type ProcessaMotorIA interface {
	SendText(text, apikey string) (string, error)
}

type ValidadorMensagem interface {
	DeveEnviarMensagem(mensagem string) bool
	SetPalavrasProibidas(palavras []string)
	ApplyFilterByRegex(mensagem string) string
	SetFilterByRegex(pattern string) error
}
