package compose

import (
	"context"

	"github.com/RMS-SH/BackgroundGo/internal/entities"
	infra "github.com/RMS-SH/BackgroundGo/internal/infra/db"
	infra_flowise "github.com/RMS-SH/BackgroundGo/internal/infra/flowise"
	"github.com/RMS-SH/BackgroundGo/internal/infra/uchat"
	"github.com/RMS-SH/BackgroundGo/internal/repositories"
	"github.com/RMS-SH/BackgroundGo/internal/usecase"
	"go.mongodb.org/mongo-driver/mongo"
)

func BackgroundCompose(
	Data []entities.MessageItem,
	apiKey string,
	db *mongo.Client,
	baseUrlUchat string,
	ctx context.Context,
) error {
	cfg := entities.NewConfig(Data[0], apiKey, baseUrlUchat)
	dbClient := infra.NewClientMongoDB(ctx, cfg, db)
	internal := uchat.NewClientUchat(ctx, cfg)
	ia := infra_flowise.NewClientFlowise(ctx, cfg)
	rp := repositories.NewProcessRepository(dbClient, internal, ctx, cfg)
	uc := usecase.NewBackgroud(rp, internal, dbClient, ctx, ia)

	return uc.ProcessaBackground(entities.Dados{Body: Data})
}
