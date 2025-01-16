package infra

import (
	"context"
	"fmt"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"

	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
)

type ClientFirebase struct {
	ctx    context.Context
	client *db.Client
}

// NewClientFirebase retorna um novo cliente para o Realtime Database.
func NewClientFirebase(ctx context.Context, app *firebase.App) (*ClientFirebase, error) {
	// Inicializa o cliente do Realtime Database a partir do App
	dbClient, err := app.Database(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar o cliente do Realtime Database: %v", err)
	}

	return &ClientFirebase{
		ctx:    ctx,
		client: dbClient,
	}, nil
}

// ArmazenaReferenciasDeArquivos faz o papel de "upsert" no Realtime Database.
func (c *ClientFirebase) ArmazenaReferenciasDeArquivos(texto, url string) error {
	ref := c.client.NewRef("referenciasTemporarias").Child(url)

	// Dados que serão gravados/atualizados
	data := map[string]interface{}{
		"texto": texto,
		"url":   url,
	}

	// Em Realtime Database, o "merge" é basicamente um Update.
	if err := ref.Update(c.ctx, data); err != nil {
		return fmt.Errorf("erro ao gravar/atualizar a referência de arquivo: %v", err)
	}

	return nil
}

// ConsultaURLReferencia lê o nó com chave = URL em "referenciasTemporarias".
func (c *ClientFirebase) ConsultaURLReferencia(url string) (string, error) {
	ref := c.client.NewRef("referenciasTemporarias").Child(url)

	// Estrutura temporária para leitura
	var result struct {
		Texto string `json:"texto"`
	}
	if err := ref.Get(c.ctx, &result); err != nil {
		return "", fmt.Errorf("erro ao consultar referência: %v", err)
	}

	if result.Texto == "" {
		return "", fmt.Errorf("nenhuma referência encontrada para a URL: %s", url)
	}

	return result.Texto, nil
}

// ConsultaDadosEmpresa busca, no nó "clientesRms", o registro cujo "workSpaceId" seja igual a workSpaceID.
func (c *ClientFirebase) ConsultaDadosEmpresa(workSpaceID string) (*entities_db.Empresa, error) {
	// Referência para a raiz de "clientesRms"
	ref := c.client.NewRef("clientesRms")

	// Mapeamento para ler múltiplos registros
	var data map[string]entities_db.Empresa

	// Faz a "query": WHERE workSpaceId = workSpaceID e LIMIT 1
	query := ref.OrderByChild("workSpaceId").EqualTo(workSpaceID).LimitToFirst(1)
	if err := query.Get(c.ctx, &data); err != nil {
		return nil, fmt.Errorf("erro ao consultar dados da empresa: %v", err)
	}

	// Se não encontrar nenhum item
	if len(data) == 0 {
		return nil, fmt.Errorf("nenhuma empresa encontrada com o ID: %s", workSpaceID)
	}

	// Pega o primeiro item do map
	for _, empresa := range data {
		return &empresa, nil
	}

	return nil, fmt.Errorf("nenhuma empresa encontrada com o ID: %s", workSpaceID)
}

// contagemConsumo é uma estrutura auxiliar para mapear o documento de "consumo".
type contagemConsumo struct {
	Workspace    string    `json:"workspace"`
	Data         time.Time `json:"data"`
	Respostas    int       `json:"respostas"`
	ImagensPdf   int       `json:"imagensPdf"`
	MinutosAudio float64   `json:"minutosAudio"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// getOrCreateConsumo localiza (pelo workspace + data) ou cria um registro base de consumo.
func (c *ClientFirebase) getOrCreateConsumo(workSpaceID string, today time.Time) (*db.Ref, error) {
	// Referência para "consumo"
	ref := c.client.NewRef("consumo")

	// Vamos filtrar por workspace e data
	var registros map[string]contagemConsumo

	query := ref.OrderByChild("workspace").EqualTo(workSpaceID)
	if err := query.Get(c.ctx, &registros); err != nil {
		return nil, fmt.Errorf("erro ao buscar documentos de consumo: %v", err)
	}

	// Varre todos os registros retornados para achar aquele cujo Data == today
	for key, val := range registros {
		// Considerando que "val.Data" foi salvo em UTC e truncado
		if val.Data.Equal(today) {
			return ref.Child(key), nil
		}
	}

	// Se não achou, cria um novo
	newConsumo := contagemConsumo{
		Workspace:    workSpaceID,
		Data:         today,
		Respostas:    0,
		ImagensPdf:   0,
		MinutosAudio: 0,
		AtualizadoEm: time.Now().UTC(),
	}

	// No Realtime Database, usamos Push() para gerar um ID único.
	newRef, err := ref.Push(c.ctx, newConsumo)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar documento de consumo: %v", err)
	}

	return newRef, nil
}

func (c *ClientFirebase) ContagemDeRespostas(workSpaceID string) error {
	// Exemplo simples; omiti parte do código do getOrCreateConsumo
	consumoRef, err := c.getOrCreateConsumo(workSpaceID, time.Now().UTC().Truncate(24*time.Hour))
	if err != nil {
		return err
	}

	// Aqui está o pulo do gato: a função de Transaction deve ter a assinatura (TransactionNode) (interface{}, error)
	return consumoRef.Transaction(c.ctx, func(tn db.TransactionNode) (interface{}, error) {
		// 1. Lê o que existe em tn
		var cur map[string]interface{}
		if err := tn.Unmarshal(&cur); err != nil {
			return nil, fmt.Errorf("erro ao fazer unmarshal: %v", err)
		}

		// Se não existir nada, cur será nil
		if cur == nil {
			cur = map[string]interface{}{}
		}

		// 2. Faz o incremento
		if val, ok := cur["respostas"].(float64); ok {
			cur["respostas"] = val + 1
		} else {
			cur["respostas"] = 1
		}

		// Só para exemplo, vamos atualizar um campo "atualizadoEm"
		cur["atualizadoEm"] = time.Now().UTC().Format(time.RFC3339)

		// 3. Retorna o map final
		return cur, nil
	})
}

// ContagemDeImagensProcessadas incrementa o campo "imagensPdf".
func (c *ClientFirebase) ContagemDeImagensProcessadas(workSpaceID string) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	consumoRef, err := c.getOrCreateConsumo(workSpaceID, today)
	if err != nil {
		return err
	}

	// Aqui está o pulo do gato: a função de Transaction deve ter a assinatura (TransactionNode) (interface{}, error)
	return consumoRef.Transaction(c.ctx, func(tn db.TransactionNode) (interface{}, error) {
		// 1. Lê o que existe em tn
		var cur map[string]interface{}
		if err := tn.Unmarshal(&cur); err != nil {
			return nil, fmt.Errorf("erro ao fazer unmarshal: %v", err)
		}

		// Se não existir nada, cur será nil
		if cur == nil {
			cur = map[string]interface{}{}
		}

		// 2. Faz o incremento
		if val, ok := cur["respostas"].(float64); ok {
			cur["imagensPdf"] = val + 1
		} else {
			cur["imagensPdf"] = 1
		}

		// Só para exemplo, vamos atualizar um campo "atualizadoEm"
		cur["atualizadoEm"] = time.Now().UTC().Format(time.RFC3339)

		// 3. Retorna o map final
		return cur, nil
	})
}

// ContagemDeArquivosProcessados também incrementa o campo "imagensPdf".
func (c *ClientFirebase) ContagemDeArquivosProcessados(workSpaceID string) error {
	// Como a lógica é a mesma, reutilizamos a função acima
	return c.ContagemDeImagensProcessadas(workSpaceID)
}

func (c *ClientFirebase) ContagemDeMinutosAudio(minutos float64, workSpaceID string) error {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	consumoRef, err := c.getOrCreateConsumo(workSpaceID, today)
	if err != nil {
		return err
	}

	// Aqui está o pulo do gato: a função de Transaction deve ter a assinatura (TransactionNode) (interface{}, error)
	return consumoRef.Transaction(c.ctx, func(tn db.TransactionNode) (interface{}, error) {
		// 1. Lê o que existe em tn
		var cur map[string]interface{}
		if err := tn.Unmarshal(&cur); err != nil {
			return nil, fmt.Errorf("erro ao fazer unmarshal: %v", err)
		}

		// Se não existir nada, cur será nil
		if cur == nil {
			cur = map[string]interface{}{}
		}

		// 2. Faz o incremento
		if val, ok := cur["respostas"].(float64); ok {
			cur["imagensPdf"] = val + minutos
		} else {
			cur["imagensPdf"] = minutos
		}

		// Só para exemplo, vamos atualizar um campo "atualizadoEm"
		cur["atualizadoEm"] = time.Now().UTC().Format(time.RFC3339)

		// 3. Retorna o map final
		return cur, nil
	})
}
