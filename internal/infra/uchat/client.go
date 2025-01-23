package uchat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	ad "github.com/RMS-SH/BackgroundGo/adapters"

	dto "github.com/RMS-SH/BackgroundGo/internal/dto/uchat"
	"github.com/RMS-SH/BackgroundGo/internal/entities"
)

type ClientUchat struct {
	ctx context.Context
	cfg entities.Config
}

func NewClientUchat(ctx context.Context, cfg entities.Config) *ClientUchat {
	return &ClientUchat{ctx: ctx, cfg: cfg}
}

func (c *ClientUchat) EnviaMensagem(message string) error {
	// Monta o corpo da requisição
	type RequestBody struct {
		UserNS  string `json:"user_ns"`
		Trigger string `json:"field_name"`
		Content string `json:"value"`
	}

	requestBody := &RequestBody{
		UserNS:  c.cfg.UUIDUser,
		Trigger: "entrega_mensagem_ia",
		Content: message,
	}

	// Serializa o corpo para JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("erro ao serializar o corpo da requisição: %v", err)
	}

	// Define a URL
	url := fmt.Sprintf("%s/api/subscriber/set-user-field-by-name", c.cfg.BaseUrlUchat)

	// Cria a requisição HTTP
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("erro ao criar a requisição HTTP: %v", err)
	}

	// Define os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.ApiKeyBot))

	// Opcional: Configurar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Envia a requisição
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar a requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("erro ao ler a resposta da requisição: %v", err)
	}

	// Verifica o status da resposta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("requisição falhou com status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil

}

func (c *ClientUchat) ConsultaReferenciaConversa(referenceId string) (*ad.ReferenceReturn, error) {
	// Monta o corpo da requisição
	type RequestBody struct {
		Mids []string `json:"mids"`
	}

	requestBody := &RequestBody{
		Mids: []string{referenceId},
	}

	// Serializa o corpo para JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar o corpo da requisição: %v", err)
	}

	// Define a URL
	url := fmt.Sprintf("%s/api/subscriber/chat-messages-by-mids", c.cfg.BaseUrlUchat)

	// Cria a requisição HTTP
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("erro ao criar a requisição HTTP: %v", err)
	}

	// Define os cabeçalhos
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.cfg.ApiKeyBot))

	// Opcional: Configurar cliente HTTP com timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Envia a requisição
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao enviar a requisição HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler a resposta da requisição: %v", err)
	}

	response := dto.ReponseUchatMid{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, fmt.Errorf("erro ao deserializar a resposta: %v", err)
	}

	// Verifica o status da resposta
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("requisição falhou com status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var Return *ad.ReferenceReturn
	if response.Data[0].MsgType != "text" {
		Return = ad.NewAdapterReferenceReturn(
			response.Data[0].MsgType,
			response.Data[0].Payload.URL,
		)
	} else {
		Return = ad.NewAdapterReferenceReturn(
			response.Data[0].MsgType,
			response.Data[0].Payload.Text,
		)
	}

	return Return, nil

}
