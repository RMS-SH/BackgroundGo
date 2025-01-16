package entities_db

type Empresa struct {
	Nome                 string            `json:"nome"`
	SkOpenAI             string            `json:"skOpenAi"`
	QuantidadeRespostas  int32             `json:"quantidadeRespostas"`
	QuantidadeMinutos    float64           `json:"quantidadeMinutos"`
	QuantidadeImagensPDF int32             `json:"quantidadeImagensPdf"`
	BaseUrlUchat         string            `json:"baseUrlUchat"`
	WorkSpaceID          string            `json:"workSpaceId"`
	Vars                 map[string]string `json:"vars"`
}
