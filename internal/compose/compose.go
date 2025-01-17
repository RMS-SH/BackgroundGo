package compose

import (
	"context"

	firebase "firebase.google.com/go"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	infra "github.com/RMS-SH/BackgroundGo/internal/infra/db"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	infra_flowise "github.com/RMS-SH/BackgroundGo/internal/infra/flowise"
	"github.com/RMS-SH/BackgroundGo/internal/infra/uchat"
	"github.com/RMS-SH/BackgroundGo/internal/repositories"
	"github.com/RMS-SH/BackgroundGo/internal/usecase"
)

func BackgroundCompose(
	Data []entities.MessageItem,
	apiKey string,
	db *firebase.App,
	baseUrlUchat string,
	ctx context.Context,
	dadosCliente entities_db.Empresa,
) error {
	cfg := entities.NewConfig(Data[0], apiKey, baseUrlUchat)
	dbClient, _ := infra.NewClientFirebase(ctx, db)
	internal := uchat.NewClientUchat(ctx, cfg)
	ia := infra_flowise.NewClientFlowise(ctx, cfg)
	rp := repositories.NewProcessRepository(dbClient, internal, ctx, cfg)
	uc := usecase.NewBackgroud(rp, internal, dbClient, ctx, ia, dadosCliente)

	return uc.ProcessaBackground(entities.Dados{Body: Data})
}

func ConsultaDadosEmpresaCompose(
	db *firebase.App,
	ctx context.Context,
	workSpaceID string,
) (*entities_db.Empresa, error) {
	dbClient, _ := infra.NewClientFirebase(ctx, db)
	return dbClient.ConsultaDadosEmpresa(workSpaceID)
}
