package usecase

import (
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	"github.com/RMS-SH/BackgroundGo/internal/interfaces"
)

type ConsultasDB struct {
	DB interfaces.DB
}

func NewConsultasDB(db interfaces.DB) *ConsultasDB {
	return &ConsultasDB{DB: db}
}

func (b *ConsultasDB) ConsultaDadosEmpresa(workSpaceID string) (*entities_db.Empresa, error) {

	resp, err := b.DB.ConsultaDadosEmpresa(workSpaceID)
	if err != nil {
		return nil, err
	}

	return resp, nil

}
