package dto_rms

import (
	"time"
)

// Struct para representar o formato solicitado
type RequestRMS struct {
	Body struct {
		Question  string            `json:"question"`
		Vars      map[string]string `json:"vars"`
		SessionID string            `json:"sessionId"`
		OpenAIKey string            `json:"openaiKey"`
	} `json:"Body"`
	UserNs string `json:"UserNs"`
	URLRms string `json:"https://rms-836463411403.us-east1.run.app"`
}

// Função para criar e retornar a struct no formato desejado
func CreateRequest(
	nome string,
	telefone string,
	UUIDUser string,
	question string,
	url string,
	extraVars ...map[string]string,

) RequestRMS {
	response := RequestRMS{}

	// Preenchendo os dados básicos
	response.Body.Question = question
	response.Body.SessionID = UUIDUser

	// Inicializando as variáveis padrão
	vars := map[string]string{
		"nome":            nome,                                     // Obtém do extraVars
		"telefone":        telefone,                                 // Obtém do extraVars
		"user_ns":         UUIDUser,                                 // Obtém do extraVars
		"data_hora_atual": time.Now().Format("02/01/2006 15:04:05"), // Formato DD/MM/YYYY HH:MM:SS
		"saudacao":        getSaudacao(time.Now().Hour()),           // Saudação baseada na hora
		"dia_semana":      getDiaSemana(time.Now().Weekday()),       // Dia da semana em português
	}

	// Se extraVars foi fornecido, adiciona-as às Vars
	if len(extraVars) > 0 && extraVars[0] != nil {
		for key, value := range extraVars[0] {
			vars[key] = value
		}
	}

	response.Body.Vars = vars

	// Atribuindo UserNs e URLFlowise a partir de extraVars ou deixando vazio se não estiver presente
	response.UserNs = UUIDUser
	response.URLRms = url

	return response
}

// Função auxiliar para determinar a saudação com base na hora do dia
func getSaudacao(hour int) string {
	switch {
	case hour >= 5 && hour < 12:
		return "Bom dia"
	case hour >= 12 && hour < 18:
		return "Boa tarde"
	case hour >= 18 && hour < 22:
		return "Boa noite"
	default:
		return "Olá"
	}
}

// Função auxiliar para obter o dia da semana em português
func getDiaSemana(day time.Weekday) string {
	dias := map[time.Weekday]string{
		time.Sunday:    "domingo",
		time.Monday:    "segunda-feira",
		time.Tuesday:   "terça-feira",
		time.Wednesday: "quarta-feira",
		time.Thursday:  "quinta-feira",
		time.Friday:    "sexta-feira",
		time.Saturday:  "sábado",
	}
	return dias[day]
}
