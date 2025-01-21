package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dto_rms "github.com/RMS-SH/BackgroundGo/internal/dto/rms"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
)

type ClientRMS struct {
	ctx context.Context
	cfg entities.Config
}

func NewClientRMSAI(ctx context.Context, cfg entities.Config) *ClientRMS {
	return &ClientRMS{ctx: ctx, cfg: cfg}
}

// ResponseItem representa cada objeto dentro do array de resposta.
type ResponseItem struct {
	Response string `json:"response"`
}

func (c *ClientRMS) SendText(text string) (string, error) {
	// Criando o request usando o dto_RMSAI.CreateRequest
	body := dto_rms.CreateRequest(
		c.cfg.Nome,         // nome
		c.cfg.Telefone,     // telefone
		c.cfg.UUIDUser,     // UUIDUser
		text,               // question
		c.cfg.URLMotorIA,   // url
		c.cfg.ApiKeyOpenIA, // openaiKey
		c.cfg.ExtraVars,    // extraVars
	)

	// Convertendo o body para JSON
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("erro ao converter body para JSON: %v", err)
	}

	// Fazendo o request HTTP
	req, err := http.NewRequestWithContext(c.ctx, "POST", "http://localhost:8080/filallm", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("erro ao criar request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.BearerTokenRMS))

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

	return Response.Response, nil
}
