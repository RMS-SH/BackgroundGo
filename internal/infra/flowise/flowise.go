package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dto_flowise "github.com/RMS-SH/BackgroundGo/internal/dto/flowise"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
)

type ClientFlowise struct {
	ctx context.Context
	cfg entities.Config
}

func NewClientFlowise(ctx context.Context, cfg entities.Config) *ClientFlowise {
	return &ClientFlowise{ctx: ctx, cfg: cfg}
}

// ResponseItem representa cada objeto dentro do array de resposta.
type ResponseItem struct {
	AgentReasoning []AgentReasoning `json:"agentReasoning"`
	ChatID         string           `json:"chatId"`
	ChatMessageID  string           `json:"chatMessageId"`
	MemoryType     string           `json:"memoryType"`
	Question       string           `json:"question"`
	SessionID      string           `json:"sessionId"`
	Text           string           `json:"text"`
}

// AgentReasoning representa cada objeto dentro do array "agentReasoning".
type AgentReasoning struct {
	AgentName       string            `json:"agentName"`
	Artifacts       []interface{}     `json:"artifacts"` // Pode ajustar o tipo conforme necessário
	Messages        []string          `json:"messages"`
	NodeID          string            `json:"nodeId"`
	NodeName        string            `json:"nodeName"`
	SourceDocuments []interface{}     `json:"sourceDocuments"` // Pode ajustar o tipo conforme necessário
	State           map[string]string `json:"state"`
	UsedTools       []interface{}     `json:"usedTools"` // Pode ajustar o tipo conforme necessário
}

func (c *ClientFlowise) SendText(text string) (string, error) {
	// Criando o request usando o dto_flowise.CreateRequest
	body := dto_flowise.CreateRequest(
		c.cfg.Nome,       // nome
		c.cfg.Telefone,   // telefone (vazio pois não está disponível no cfg)
		c.cfg.UUIDUser,   // UUIDUser
		text,             // question
		c.cfg.URLMotorIA, // url
	)

	// Convertendo o body para JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("erro ao converter body para JSON: %v", err)
	}

	// Fazendo o request HTTP
	req, err := http.NewRequestWithContext(c.ctx, "POST", "http://localhost:8080/filaflowise", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer v_RI9r7-+)hMoedU@~H[jUGJctDhU;")

	// Executando o request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("erro ao executar request: %v", err)
	}
	defer resp.Body.Close()

	// Verificando o status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request falhou com status code: %d", resp.StatusCode)
	}

	var Response ResponseItem
	err = json.NewDecoder(resp.Body).Decode(&Response)
	if err != nil {
		return "", fmt.Errorf("erro ao decodificar resposta: %v", err)
	}

	return Response.Text, nil
}
