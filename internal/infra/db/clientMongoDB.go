package infra

import (
	"context"
	"fmt"
	"time"

	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ClientMongoDB struct {
	ctx context.Context
	db  *mongo.Client
}

func NewClientMongoDB(ctx context.Context, db *mongo.Client) *ClientMongoDB {
	return &ClientMongoDB{ctx: ctx, db: db}
}

func (c *ClientMongoDB) ArmazenaReferenciasDeArquivos(texto, url string) error {
	// Define a coleção
	collection := c.db.Database("rms").Collection("referencias_temporarias")

	// Define o filtro e o documento de atualização
	filter := struct {
		URL string `bson:"url"`
	}{
		URL: url,
	}

	update := struct {
		Set struct {
			Texto string `bson:"texto"`
			URL   string `bson:"url"`
		} `bson:"$set"`
	}{
		Set: struct {
			Texto string `bson:"texto"`
			URL   string `bson:"url"`
		}{
			Texto: texto,
			URL:   url,
		},
	}

	// Realiza o upsert
	_, err := collection.UpdateOne(c.ctx, filter, update, options.Update().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientMongoDB) ConsultaURLReferencia(url string) (string, error) {
	// Define a coleção
	collection := c.db.Database("rms").Collection("referencias_temporarias")

	// Define o filtro para a consulta
	filter := struct {
		URL string `bson:"url"`
	}{
		URL: url,
	}

	// Estrutura para armazenar o resultado
	var result struct {
		Texto string `bson:"texto"`
	}

	// Realiza a consulta
	err := collection.FindOne(c.ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Retorna erro se não encontrar documentos
			return "", fmt.Errorf("nenhuma referência encontrada para a URL: %s", url)
		}
		// Retorna outros erros
		return "", err
	}

	// Retorna o texto encontrado
	return result.Texto, nil
}

func (c *ClientMongoDB) ConsultaDadosEmpresa(workSpaceID string) (*entities_db.Empresa, error) {
	// Define a coleção
	collection := c.db.Database("rms").Collection("clientesRms")

	// Define o filtro para a consulta
	filter := bson.M{
		"workSpaceId": workSpaceID, // Supondo que o campo _id seja uma string. Ajuste se for ObjectID.
	}

	// Estrutura para armazenar o resultado
	var result entities_db.Empresa

	// Realiza a consulta
	err := collection.FindOne(c.ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("nenhuma empresa encontrada com o ID: %s", workSpaceID)
		}
		return nil, err
	}

	// Retorna o nome da empresa encontrada
	return &result, nil
}

func (c *ClientMongoDB) ContagemDeRespostas(workSpaceID string) error {
	// Define a coleção
	collection := c.db.Database("rms").Collection("consumo")

	// Obtém a data atual sem a hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Define o filtro para encontrar o documento do workspace e da data atual
	filter := bson.M{
		"workspace": workSpaceID,
		"data":      today,
	}

	// Define a atualização para incrementar o campo 'respostas' e atualizar 'atualizado_em'
	update := bson.M{
		"$inc": bson.M{
			"respostas": 1,
		},
		"$set": bson.M{
			"atualizado_em": time.Now().UTC(),
		},
	}

	// Define as opções para upsert
	opts := options.Update().SetUpsert(true)

	// Realiza a operação de atualização
	_, err := collection.UpdateOne(c.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao contar respostas: %v", err)
	}

	return nil
}

func (c *ClientMongoDB) ContagemDeMinutosAudio(minutos float64, workSpaceID string) error {
	// Define a coleção
	collection := c.db.Database("rms").Collection("consumo")

	// Obtém a data atual sem a hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Define o filtro para encontrar o documento do workspace e da data atual
	filter := bson.M{
		"workspace": workSpaceID,
		"data":      today,
	}

	// Define a atualização para incrementar o campo 'minutos_audio' e atualizar 'atualizado_em'
	update := bson.M{
		"$inc": bson.M{
			"minutos_audio": minutos,
		},
		"$set": bson.M{
			"atualizado_em": time.Now().UTC(),
		},
	}

	// Define as opções para upsert
	opts := options.Update().SetUpsert(true)

	// Realiza a operação de atualização
	_, err := collection.UpdateOne(c.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao contar minutos de áudio: %v", err)
	}

	return nil
}

func (c *ClientMongoDB) ContagemDeImagensProcessadas(workSpaceID string) error {
	// Define a coleção
	collection := c.db.Database("rms").Collection("consumo")

	// Obtém a data atual sem a hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Define o filtro para encontrar o documento do workspace e da data atual
	filter := bson.M{
		"workspace": workSpaceID,
		"data":      today,
	}

	// Define a atualização para incrementar o campo 'imagens' e atualizar 'atualizado_em'
	update := bson.M{
		"$inc": bson.M{
			"imagens": 1,
		},
		"$set": bson.M{
			"atualizado_em": time.Now().UTC(),
		},
	}

	// Define as opções para upsert
	opts := options.Update().SetUpsert(true)

	// Realiza a operação de atualização
	_, err := collection.UpdateOne(c.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao contar imagens processadas: %v", err)
	}

	return nil
}

func (c *ClientMongoDB) ContagemDeArquivosProcessados(workSpaceID string) error {
	// Define a coleção
	collection := c.db.Database("rms").Collection("consumo")

	// Obtém a data atual sem a hora
	today := time.Now().UTC().Truncate(24 * time.Hour)

	// Define o filtro para encontrar o documento do workspace e da data atual
	filter := bson.M{
		"workspace": workSpaceID,
		"data":      today,
	}

	// Define a atualização para incrementar o campo 'pdf' e atualizar 'atualizado_em'
	update := bson.M{
		"$inc": bson.M{
			"pdf": 1,
		},
		"$set": bson.M{
			"atualizado_em": time.Now().UTC(),
		},
	}

	// Define as opções para upsert
	opts := options.Update().SetUpsert(true)

	// Realiza a operação de atualização
	_, err := collection.UpdateOne(c.ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("erro ao contar arquivos processados: %v", err)
	}

	return nil
}
