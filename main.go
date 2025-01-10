package backgroundgo

import (
	"context"

	"github.com/RMS-SH/BackgroundGo/internal/compose"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

func Background(
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
