package backgroundgo

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/RMS-SH/BackgroundGo/internal/compose"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
)

type Backgroud struct {
	db *firestore.Client
}

func NewBackgroud(db *firestore.Client) *Backgroud {
	return &Backgroud{db: db}
}

func (bk *Backgroud) Proccess(
	Data []entities.MessageItem,
	apiKey string,
	db *firestore.Client,
	baseUrlUchat string,
	ctx context.Context,
	dadosCliente entities_db.Empresa,
	extraVars map[string]string,
	tokenRMS string,
) error {
	return compose.BackgroundCompose(
		Data,
		apiKey,
		db,
		baseUrlUchat,
		ctx,
		dadosCliente,
		extraVars,
		tokenRMS,
	)
}

func (bk *Backgroud) ConsultaDadosWorkSpaceID(dados []entities.MessageItem, ctx context.Context) (*entities_db.Empresa, error) {
	DadosEmpresa, err := compose.ConsultaDadosEmpresaCompose(bk.db, ctx, dados[0].IDWorkSpace)
	if err != nil {
		return nil, err
	}

	return DadosEmpresa, nil
}

func (bk *Backgroud) GerarEntiteInterna(
	typeValue string,
	content string,
	reference string,
	nome string,
	telefone string,
	userNS string,
	motorIA string,
	bearerMotorIA string,
	dataHoraAtual string,
	diaSemana string,
	url string,
	urlError string,
	urlBackground string,
	jsonString string,
	nomeWorkspace string,
	idWorkspace string,
	apiKeyBot string,
	respondeAudio bool,
	apiKeyElevenLabs string,
	voiceIdElevenLabs string,
	promptArquivo string,
	promptImagem string,
	arraysVarRetorno string,
) entities.MessageItem {
	return entities.MessageItem{
		Type:              typeValue,
		Content:           content,
		Reference:         reference,
		Nome:              nome,
		Telefone:          telefone,
		UserNS:            userNS,
		MotorIA:           motorIA,
		BearerMotorIA:     bearerMotorIA,
		DataHoraAtual:     dataHoraAtual,
		DiaSemana:         diaSemana,
		URL:               url,
		URLError:          urlError,
		URLBackGround:     urlBackground,
		JSONString:        jsonString,
		NomeWorkspace:     nomeWorkspace,
		IDWorkSpace:       idWorkspace,
		ApiKeyBot:         apiKeyBot,
		RespondeAudio:     respondeAudio,
		ApiKeyElevenLabs:  apiKeyElevenLabs,
		VoiceIdElevenLabs: voiceIdElevenLabs,
		PromptArquivo:     promptArquivo,
		PromptImagem:      promptImagem,
		ArraysVarRetorno:  arraysVarRetorno,
	}
}

func (bk *Backgroud) RetornaStructInternal() []entities.MessageItem {
	return []entities.MessageItem{}
}
