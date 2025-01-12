package compose

import (
	"context"

	"github.com/RMS-SH/BackgroundGo/internal/entities"
	infra "github.com/RMS-SH/BackgroundGo/internal/infra/db"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	infra_flowise "github.com/RMS-SH/BackgroundGo/internal/infra/flowise"
	"github.com/RMS-SH/BackgroundGo/internal/infra/uchat"
	"github.com/RMS-SH/BackgroundGo/internal/repositories"
	"github.com/RMS-SH/BackgroundGo/internal/usecase"
	"github.com/RMS-SH/BackgroundGo/internal/validators"
	"go.mongodb.org/mongo-driver/mongo"
)

func BackgroundCompose(
	Data []entities.MessageItem,
	apiKey string,
	db *mongo.Client,
	baseUrlUchat string,
	ctx context.Context,
	dadosCliente entities_db.Empresa,
) error {
	cfg := entities.NewConfig(Data[0], apiKey, baseUrlUchat)
	dbClient := infra.NewClientMongoDB(ctx, db)
	internal := uchat.NewClientUchat(ctx, cfg)
	ia := infra_flowise.NewClientFlowise(ctx, cfg)
	validador := validators.NewMessageValidator()
	rp := repositories.NewProcessRepository(dbClient, internal, ctx, cfg)
	uc := usecase.NewBackgroud(rp, internal, dbClient, ctx, ia, validador, dadosCliente)

	return uc.ProcessaBackground(entities.Dados{Body: Data})
}

func ConsultaDadosEmpresaCompose(
	db *mongo.Client,
	ctx context.Context,
	workSpaceID string,
) (*entities_db.Empresa, error) {
	dbClient := infra.NewClientMongoDB(ctx, db)
	return dbClient.ConsultaDadosEmpresa(workSpaceID)
}
