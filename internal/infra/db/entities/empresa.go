package entities_db

type Empresa struct {
	Nome                 string  `bson:"nome"`
	SkOpenAI             string  `bson:"skOpenAI"`
	QuantidadeRespostas  int32   `bson:"quantidadeRespostas"`
	QuantidadeMinutos    float64 `bson:"quantidadeMinutos"`
	QuantidadeImagensPDF int32   `bson:"quantidadeImagensPdf"`
	WorkSpaceID          string  `bson:"workSpaceId"`
}
