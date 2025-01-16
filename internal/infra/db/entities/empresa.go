package entities_db

type Empresa struct {
	Nome                 string   `firestore:"nome"`
	SkOpenAI             string   `firestore:"skOpenAi"`
	QuantidadeRespostas  int32    `firestore:"quantidadeRespostas"`
	QuantidadeMinutos    float64  `firestore:"quantidadeMinutos"`
	QuantidadeImagensPDF int32    `firestore:"quantidadeImagensPdf"`
	BaseUrlUchat         string   `firestore:"baseUrlUchat"`
	WorkSpaceID          string   `firestore:"workSpaceId"`
	PalavrasProibidas    []string `firestore:"palavrasProibidas"`
	PalavraRGEX          string   `firestore:"palavraRGEX"`
}
