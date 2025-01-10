package backgroundgo

import (
	"context"

	"github.com/RMS-SH/BackgroundGo/internal/compose"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

type Backgroud struct {
	db *mongo.Client
}

func NewBackgroud(db *mongo.Client) *Backgroud {
	return &Backgroud{db: db}
}

func (bk *Backgroud) Proccess(
	Data []entities.MessageItem,
	apiKey string,
	db *mongo.Client,
	baseUrlUchat string,
	ctx context.Context,
) error {
	return compose.BackgroundCompose(
		Data,
		apiKey,
		db,
		baseUrlUchat,
		ctx,
	)
}

func (bk *Backgroud) ConsultaDadosWorkSpaceID(dados []entities.MessageItem, ctx context.Context) (*entities_db.Empresa, error) {
	DadosEmpresa, err := compose.ConsultaDadosEmpresaCompose(bk.db, ctx, dados[0].IDWorkSpace)
	if err != nil {
		return nil, err
	}

	return DadosEmpresa, nil
}
