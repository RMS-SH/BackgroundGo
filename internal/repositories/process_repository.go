package repositories

import (
	"context"
	"fmt"

	"github.com/RMS-SH/BackgroundGo/internal/entities"
	"github.com/RMS-SH/BackgroundGo/internal/interfaces"
	openia "github.com/RMS-SH/OpenIA"
)

type ProcessRepository struct {
	db       interfaces.DB
	internal interfaces.Internal
	ctx      context.Context
	cfg      entities.Config
}

func NewProcessRepository(db interfaces.DB, internal interfaces.Internal, ctx context.Context, cfg entities.Config) *ProcessRepository {
	return &ProcessRepository{
		db:       db,
		internal: internal,
		ctx:      ctx,
		cfg:      cfg,
	}
}

func (r *ProcessRepository) ProcessaTexto(msg entities.MessageItem) (string, error) {
	// Função auxiliar para verificar tipos que precisam de URL
	precisaConsultarURL := func(tipo string) bool {
		tiposComURL := map[string]bool{
			"audio": true,
			"image": true,
			"file":  true,
		}
		return tiposComURL[tipo]
	}

	// Se não tem referência apenas retorna o Texto
	if msg.Reference == "" {
		return msg.Content, nil
	}

	// Caso tenha referência precisamos ver se consguimos consultar! pois as referências ficam gravadas pro 30 dias!
	referencia, err := r.internal.ConsultaReferenciaConversa(msg.Reference)
	if err != nil {
		return msg.Content, nil // retorna apenas o conteúdo original se houver erro
	}

	// Após consultar a Referência o Metodo ConsultaReferenciaConversa ele retorna Tipo e Value
	// Caso a referência ele seja do tipo Texto não se faz necessário consultar no banco de dados.
	if !precisaConsultarURL(referencia.Tipo) {
		return fmt.Sprintf("Contexto:%s\n", referencia.Value) + "Pergunta : " + msg.Content, nil
	}

	// Consulta URL apenas para tipos específicos, se é uma referência já foi teoricamente processado em outro momento.
	// ENtão as URL são armazenadas no banco de dados e só precisamos consultar o banco de dados.
	ReferenciaURL, err := r.db.ConsultaURLReferencia(referencia.Value)
	if err != nil {
		return msg.Content + referencia.Value, nil
	}

	return fmt.Sprintf("Contexto:%s\n", ReferenciaURL) + "Pergunta : " + msg.Content, nil
}

func (r *ProcessRepository) ProcessaImagem(msg entities.MessageItem) (string, error) {

	ImageToText, err := openia.AnalisaImage(r.ctx, "OpenIA", msg.Content, r.cfg.ApiKeyOpenIA, msg.PromptImagem, "", "")
	if err != nil {
		return "", err
	}

	GravaURLReferencia := r.db.ArmazenaReferenciasDeArquivos(ImageToText.Text, msg.Content)
	if GravaURLReferencia != nil {
		return "", GravaURLReferencia
	}

	ContagemDeImagensProcessadas := r.db.ContagemDeImagensProcessadas()
	if ContagemDeImagensProcessadas != nil {
		return "", ContagemDeImagensProcessadas
	}

	return ImageToText.Text, nil
}

func (r *ProcessRepository) ProcessaAudio(msg entities.MessageItem) (string, error) {

	AudioTranscription, err := openia.AudioTranscription(r.ctx, "OpenIA", r.cfg.ApiKeyOpenIA, msg.Content, "", "")
	if err != nil {
		return "", err
	}

	GravaURLReferencia := r.db.ArmazenaReferenciasDeArquivos(AudioTranscription.Text, msg.Content)
	if GravaURLReferencia != nil {
		return AudioTranscription.Text, GravaURLReferencia
	}

	err = r.db.ContagemDeMinutosAudio(AudioTranscription.DurationSegundos)
	if err != nil {
		return AudioTranscription.Text, err
	}

	return AudioTranscription.Text, nil
}

func (r *ProcessRepository) ProcessaFile(texto entities.MessageItem) (string, error) {

	DetalhesDocumento, err := openia.InterpretacaoPDFAssistente(r.ctx, r.cfg.PromptArquivo, texto.Content, r.cfg.ApiKeyOpenIA)
	if err != nil {
		return "", err
	}

	GravaURLReferencia := r.db.ArmazenaReferenciasDeArquivos(DetalhesDocumento.(string), texto.Content)
	if GravaURLReferencia != nil {
		return DetalhesDocumento.(string), GravaURLReferencia
	}

	err = r.db.ContagemDeImagensProcessadas()
	if err != nil {
		return DetalhesDocumento.(string), err
	}

	return DetalhesDocumento.(string), nil
}
