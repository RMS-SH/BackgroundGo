package entities

type MessageItem struct {
	Type              string `json:"type,omitempty"`
	Content           string `json:"content,omitempty"`
	Reference         string `json:"mid,omitempty"`
	Nome              string `json:"nome,omitempty"`
	Telefone          string `json:"telefone,omitempty"`
	UserNS            string `json:"userNs,omitempty"`
	MotorIA           string `json:"motorIA,omitempty"`
	BearerMotorIA     string `json:"bearerMotorIA,omitempty"`
	DataHoraAtual     string `json:"dataHoraAtual,omitempty"`
	DiaSemana         string `json:"diaSemana,omitempty"`
	URL               string `json:"url,omitempty"`
	URLError          string `json:"urlError,omitempty"`
	URLBackGround     string `json:"urlBackground,omitempty"`
	JSONString        string `json:"jsonString,omitempty"`
	NomeWorkspace     string `json:"nomeWorkspace,omitempty"`
	IDWorkSpace       string `json:"workspaceId,omitempty"`
	ApiKeyBot         string `json:"apiKeyBot,omitempty"`
	RespondeAudio     bool   `json:"respondeAudio,omitempty"`
	ApiKeyElevenLabs  string `json:"apiKeyElevenLabs,omitempty"`
	VoiceIdElevenLabs string `json:"voiceIdElevenLabs,omitempty"`
	PromptArquivo     string `json:"promptArquivo,omitempty"`
	PromptImagem      string `json:"promptImagem,omitempty"`
	ArraysVarRetorno  string `json:"arrayVars,omitempty"`
}

type Dados struct {
	Body []MessageItem `json:"body"`
}
