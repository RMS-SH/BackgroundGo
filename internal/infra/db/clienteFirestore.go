package infra

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"

	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ClientFirestore struct {
	ctx    context.Context
	client *firestore.Client
}

// NewClientFirestore retorna um novo cliente para o Firestore.
func NewClientFirestore(ctx context.Context, client *firestore.Client) *ClientFirestore {
	return &ClientFirestore{
		ctx:    ctx,
		client: client,
	}
}

// ArmazenaReferenciasDeArquivos faz o papel do 'upsert' no Firestore.
// Cria ou atualiza o documento cuja chave será igual à URL fornecida.
func (c *ClientFirestore) ArmazenaReferenciasDeArquivos(texto, url string) error {
	// Gera um hash da URL para usar como ID do documento
	urlHash := sha256.Sum256([]byte(url))
	docID := hex.EncodeToString(urlHash[:])

	collection := c.client.Collection("referenciasTemporarias")
	docRef := collection.Doc(docID)

	data := map[string]interface{}{
		"texto":    texto,
		"url":      url,
		"criadoEm": firestore.ServerTimestamp,
	}

	_, err := docRef.Set(c.ctx, data, firestore.MergeAll)
	if err != nil {
		return fmt.Errorf("erro ao gravar/atualizar a referência de arquivo: %v", err)
	}

	return nil
}

func (c *ClientFirestore) ConsultaURLReferencia(url string) (string, error) {
	// Gera o mesmo hash da URL usado no ArmazenaReferenciasDeArquivos
	urlHash := sha256.Sum256([]byte(url))
	docID := hex.EncodeToString(urlHash[:])

	// A coleção de destino
	collection := c.client.Collection("referenciasTemporarias")

	// Documento com ID = hash da URL
	docRef := collection.Doc(docID)

	// Tenta obter o documento
	docSnap, err := docRef.Get(c.ctx)
	if err != nil {
		// Verificamos o código do erro
		if status.Code(err) == codes.NotFound {
			return "", fmt.Errorf("nenhuma referência encontrada para a URL: %s", url)
		}
		return "", fmt.Errorf("erro ao consultar referência: %v", err)
	}

	// Deserializa o resultado em uma struct ou map
	var result struct {
		Texto string `firestore:"texto"`
	}
	err = docSnap.DataTo(&result)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar dados: %v", err)
	}

	return result.Texto, nil
}

// ConsultaDadosEmpresa busca o documento na coleção "clientesRms" filtrando por "workSpaceId".
func (c *ClientFirestore) ConsultaDadosEmpresa(workSpaceID string) (*entities_db.Empresa, error) {
	// A coleção de destino
	collection := c.client.Collection("clientesRms")

	// Faz a query: WHERE workSpaceId = workSpaceID
	query := collection.Where("workSpaceId", "==", workSpaceID).Limit(1)

	// Executa a query
	iter := query.Documents(c.ctx)
	defer iter.Stop()

	// Lê o primeiro resultado
	docSnap, err := iter.Next()
	if err == iterator.Done {
		return nil, fmt.Errorf("nenhuma empresa encontrada com o ID: %s", workSpaceID)
	}
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar dados da empresa: %v", err)
	}

	// Deserializa o documento em uma struct
	var empresa entities_db.Empresa
	if err = docSnap.DataTo(&empresa); err != nil {
		return nil, fmt.Errorf("erro ao decodificar dados da empresa: %v", err)
	}

	// Retorna a empresa encontrada
	return &empresa, nil
}

// ContagemDeRespostas atualiza (ou cria) o documento na coleção "consumo"
// cujo campo "workspace" seja igual a workSpaceID e "data" seja igual ao dia de hoje,
// incrementando o campo "respostas".
func (c *ClientFirestore) ContagemDeRespostas(workSpaceID string) error {
	// A coleção de destino
	collection := c.client.Collection("consumo")

	// Data atual (zerando hora, minuto, segundo, etc.)
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Query para encontrar se já existe doc para (workSpaceID, today)
	query := collection.
		Where("workspace", "==", workSpaceID).
		Where("data", "==", today).
		Limit(1)

	iter := query.Documents(c.ctx)
	defer iter.Stop()

	docSnap, err := iter.Next()
	if err == iterator.Done {
		// Documento não existe, então criamos um novo
		newDoc := map[string]interface{}{
			"workspace":     workSpaceID,
			"data":          today,
			"respostas":     1,
			"imagensPdf":    0,
			"minutosAudio":  0,
			"atualizado_em": time.Now().UTC(),
		}
		_, _, errCreate := collection.Add(c.ctx, newDoc)
		if errCreate != nil {
			return fmt.Errorf("erro ao criar documento de contagem de respostas: %v", errCreate)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("erro ao buscar documento de contagem de respostas: %v", err)
	}

	// Documento existe; atualizamos
	docRef := collection.Doc(docSnap.Ref.ID)
	_, err = docRef.Update(c.ctx, []firestore.Update{
		{Path: "respostas", Value: firestore.Increment(1)},
		{Path: "atualizado_em", Value: time.Now().UTC()},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar contagem de respostas: %v", err)
	}

	return nil
}

// ContagemDeMinutosAudio atualiza (ou cria) o documento na coleção "consumo"
// cujo campo "workspace" seja igual a workSpaceID e "data" seja igual ao dia de hoje,
// incrementando o campo "minutosAudio".
func (c *ClientFirestore) ContagemDeMinutosAudio(minutos float64, workSpaceID string) error {
	// A coleção
	collection := c.client.Collection("consumo")

	// Data atual sem hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Query
	query := collection.
		Where("workspace", "==", workSpaceID).
		Where("data", "==", today).
		Limit(1)

	iter := query.Documents(c.ctx)
	defer iter.Stop()

	docSnap, err := iter.Next()
	if err == iterator.Done {
		// Se não existe, cria
		newDoc := map[string]interface{}{
			"workspace":     workSpaceID,
			"data":          today,
			"respostas":     1,
			"imagensPdf":    0,
			"minutosAudio":  0,
			"atualizado_em": time.Now().UTC(),
		}
		_, _, errCreate := collection.Add(c.ctx, newDoc)
		if errCreate != nil {
			return fmt.Errorf("erro ao criar documento de contagem de minutos: %v", errCreate)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("erro ao buscar documento de contagem de minutos: %v", err)
	}

	// Documento existe; atualiza
	docRef := collection.Doc(docSnap.Ref.ID)
	_, err = docRef.Update(c.ctx, []firestore.Update{
		{Path: "minutosAudio", Value: firestore.Increment(minutos)},
		{Path: "atualizadEm", Value: time.Now().UTC()},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar contagem de minutos: %v", err)
	}

	return nil
}

// ContagemDeImagensProcessadas incrementa o campo "imagensPdf".
func (c *ClientFirestore) ContagemDeImagensProcessadas(workSpaceID string) error {
	// A coleção
	collection := c.client.Collection("consumo")

	// Data atual sem hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Query
	query := collection.
		Where("workspace", "==", workSpaceID).
		Where("data", "==", today).
		Limit(1)

	iter := query.Documents(c.ctx)
	defer iter.Stop()

	docSnap, err := iter.Next()
	if err == iterator.Done {
		// Se não existe, cria
		newDoc := map[string]interface{}{
			"workspace":     workSpaceID,
			"data":          today,
			"respostas":     1,
			"imagensPdf":    0,
			"minutosAudio":  0,
			"atualizado_em": time.Now().UTC(),
		}
		_, _, errCreate := collection.Add(c.ctx, newDoc)
		if errCreate != nil {
			return fmt.Errorf("erro ao criar documento de contagem de imagens: %v", errCreate)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("erro ao buscar documento de contagem de imagens: %v", err)
	}

	// Documento existe; atualiza
	docRef := collection.Doc(docSnap.Ref.ID)
	_, err = docRef.Update(c.ctx, []firestore.Update{
		{Path: "imagensPdf", Value: firestore.Increment(1)},
		{Path: "atualizadEm", Value: time.Now().UTC()},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar contagem de imagens: %v", err)
	}

	return nil
}

// ContagemDeArquivosProcessados incrementa o campo "imagensPdf" para contabilizar PDFs ou arquivos diversos.
func (c *ClientFirestore) ContagemDeArquivosProcessados(workSpaceID string) error {
	// A coleção
	collection := c.client.Collection("consumo")

	// Data atual sem hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Query
	query := collection.
		Where("workspace", "==", workSpaceID).
		Where("data", "==", today).
		Limit(1)

	iter := query.Documents(c.ctx)
	defer iter.Stop()

	docSnap, err := iter.Next()
	if err == iterator.Done {
		// Se não existe, cria
		newDoc := map[string]interface{}{
			"workspace":     workSpaceID,
			"data":          today,
			"respostas":     1,
			"imagensPdf":    0,
			"minutosAudio":  0,
			"atualizado_em": time.Now().UTC(),
		}
		_, _, errCreate := collection.Add(c.ctx, newDoc)
		if errCreate != nil {
			return fmt.Errorf("erro ao criar documento de contagem de arquivos: %v", errCreate)
		}
		return nil
	} else if err != nil {
		return fmt.Errorf("erro ao buscar documento de contagem de arquivos: %v", err)
	}

	// Documento existe; atualiza
	docRef := collection.Doc(docSnap.Ref.ID)
	_, err = docRef.Update(c.ctx, []firestore.Update{
		{Path: "imagensPdf", Value: firestore.Increment(1)},
		{Path: "atualizadEm", Value: time.Now().UTC()},
	})
	if err != nil {
		return fmt.Errorf("erro ao atualizar contagem de arquivos: %v", err)
	}

	return nil
}
