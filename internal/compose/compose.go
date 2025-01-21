package compose

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	infra "github.com/RMS-SH/BackgroundGo/internal/infra/db"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	infra_rmsai "github.com/RMS-SH/BackgroundGo/internal/infra/rmsia"
	"github.com/RMS-SH/BackgroundGo/internal/infra/uchat"
	"github.com/RMS-SH/BackgroundGo/internal/repositories"
	"github.com/RMS-SH/BackgroundGo/internal/usecase"
	"github.com/RMS-SH/BackgroundGo/internal/validators"
)

func BackgroundCompose(
	Data []entities.MessageItem,
	apiKey string,
	db *firestore.Client,
	baseUrlUchat string,
	ctx context.Context,
	dadosCliente entities_db.Empresa,
	extraVars map[string]string,
) error {
	cfg := entities.NewConfig(Data[0], apiKey, baseUrlUchat, extraVars)
	dbClient := infra.NewClientFirestore(ctx, db)
	internal := uchat.NewClientUchat(ctx, cfg)
	ia := infra_rmsai.NewClientRMSAI(ctx, cfg)
	validador := validators.NewMessageValidator()
	rp := repositories.NewProcessRepository(dbClient, internal, ctx, cfg)
	uc := usecase.NewBackgroud(rp, internal, dbClient, ctx, ia, validador, dadosCliente)

	return uc.ProcessaBackground(entities.Dados{Body: Data})
}

func ConsultaDadosEmpresaCompose(
	db *firestore.Client,
	ctx context.Context,
	workSpaceID string,
) (*entities_db.Empresa, error) {
	dbClient := infra.NewClientFirestore(ctx, db)
	return dbClient.ConsultaDadosEmpresa(workSpaceID)
}
