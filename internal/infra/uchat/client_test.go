//go:build integration
// +build integration

package uchat

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RMS-SH/BackgroundGo/internal/entities"
	"github.com/stretchr/testify/assert"
)

// TestClientUchat_EnviaMensagem verifica o funcionamento do método EnviaMensagem
func TestClientUchat_EnviaMensagem(t *testing.T) {
	// Casos de teste
	tests := []struct {
		name           string
		message        string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
		configOverride func(cfg *entities.Config)
		clientTimeout  time.Duration
	}{
		{
			name:    "sucesso_ao_enviar_mensagem",
			message: "Olá, isso é um teste!",
			mockResponse: map[string]interface{}{
				"status": "success",
			},
			mockStatusCode: http.StatusOK,
			expectError:    false,
		},
		{
			name:    "erro_servidor",
			message: "Mensagem teste",
			mockResponse: map[string]interface{}{
				"error": "Erro interno do servidor",
			},
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
		{
			name:    "timeout_do_servidor",
			message: "Mensagem com timeout",
			mockResponse: map[string]interface{}{
				"status": "success",
			},
			mockStatusCode: http.StatusOK,
			expectError:    true,
			// Simula delay no servidor para causar timeout
			configOverride: func(cfg *entities.Config) {
				cfg.Timeout = 1 * time.Minute
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Criar servidor mock
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verificar método HTTP
				assert.Equal(t, http.MethodPost, r.Method)

				// Verificar headers
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Contains(t, r.Header.Get("Authorization"), "Bearer ")

				// Verificar corpo da requisição
				var requestBody struct {
					UserNS  string            `json:"user_ns"`
					Trigger string            `json:"trigger_name"`
					Content map[string]string `json:"context"`
				}

				decoder := json.NewDecoder(r.Body)
				err := decoder.Decode(&requestBody)
				assert.NoError(t, err)
				assert.Equal(t, "entrega_mensagem_ia", requestBody.Trigger)
				assert.Equal(t, tt.message, requestBody.Content["EntregaMensagem"])

				// Simular delay se necessário
				if tt.name == "timeout_do_servidor" {
					time.Sleep(2 * time.Second)
				}

				// Configurar resposta mock
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Configurar cliente com URL do servidor mock
			cfg := entities.Config{
				BaseUrlUchat: server.URL,
				ApiKeyBot:    "CUd5cJQoeN59tzTmZuo3zQd6j5439ySVSh8BCsbp0Seh1c6JaUP4XBTIE9ed",
				UUIDUser:     "f127583u165036325",
				Timeout:      10 * time.Second,
			}

			// Aplicar overrides se necessário
			if tt.configOverride != nil {
				tt.configOverride(&cfg)
			}

			client := NewClientUchat(context.Background(), cfg)

			// Executar teste
			err := client.EnviaMensagem(tt.message)

			// Verificar resultado
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestClientUchat_EnviaMensagemReal realiza um teste real do método EnviaMensagem sem utilizar mocks
func TestClientUchat_EnviaMensagemReal(t *testing.T) {
	// Configuração de ambiente para testes reais
	cfg := entities.Config{
		BaseUrlUchat: "https://app.chatrms.com", // URL real da API
		ApiKeyBot:    "CUd5cJQoeN59tzTmZuo3zQd6j5439ySVSh8BCsbp0Seh1c6JaUP4XBTIE9ed",
		UUIDUser:     "f127583u165036325",
		Timeout:      10 * time.Second,
	}

	client := NewClientUchat(context.Background(), cfg)

	// Mensagem de teste real
	mensagem := "Olá, esta é uma mensagem de teste real!"

	// Executar o método EnviaMensagem
	err := client.EnviaMensagem(mensagem)

	// Verificar resultado
	if err != nil {
		t.Errorf("Falha ao enviar mensagem real: %v", err)
	} else {
		t.Log("Mensagem real enviada com sucesso.")
	}
}
