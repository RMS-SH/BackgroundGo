package interfaces

import (
	ad "github.com/RMS-SH/BackgroundGo/adapters"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
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
	ConsultaDadosEmpresa() (string, error)
	ContagemDeRespostas() error
	ContagemDeMinutosAudio(minutos float64) error
	ContagemDeImagensProcessadas() error
	ContagemDeArquivosProcessados() error
}

type Internal interface {
	ConsultaReferenciaConversa(referenceId string) (*ad.ReferenceReturn, error)
}

type ProcessaMotorIA interface {
	SendText(text string) (string, error)
}
