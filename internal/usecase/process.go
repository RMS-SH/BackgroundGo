package usecase

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/RMS-SH/BackgroundGo/internal/entities"
	entities_db "github.com/RMS-SH/BackgroundGo/internal/infra/db/entities"
	"github.com/RMS-SH/BackgroundGo/internal/interfaces"
	utilitariosgorms "github.com/RMS-SH/UtilitariosGoRMS"
)

type Backgroud struct {
	Process      interfaces.Process
	Entrega      interfaces.Entrega
	DB           interfaces.DB
	ctx          context.Context
	IA           interfaces.ProcessaMotorIA
	DadosCliente entities_db.Empresa
}

func NewBackgroud(
	process interfaces.Process,
	entrega interfaces.Entrega,
	db interfaces.DB,
	ctx context.Context,
	IA interfaces.ProcessaMotorIA,
	dadosCliente entities_db.Empresa,
) *Backgroud {
	return &Backgroud{
		Process:      process,
		Entrega:      entrega,
		DB:           db,
		ctx:          ctx,
		IA:           IA,
		DadosCliente: dadosCliente,
	}
}

func (b *Backgroud) ProcessaBackground(dados entities.Dados) error {

	// Preparar slice para armazenar os resultados na ordem correta
	results := make([]string, len(dados.Body)-1)

	// Canal para captura de erros
	errCh := make(chan error, 1)

	var wg sync.WaitGroup
	wg.Add(len(dados.Body) - 1)

	for i := 1; i < len(dados.Body); i++ {
		index := i - 1 // Índice no slice de resultados
		item := dados.Body[i]

		go func(idx int, itm entities.MessageItem) {
			defer wg.Done()

			var resposta string
			var err error

			switch itm.Type {
			case "text":
				resposta, err = b.Process.ProcessaTexto(itm)
			case "image":
				resposta, err = b.Process.ProcessaImagem(itm)
			case "audio":
				resposta, err = b.Process.ProcessaAudio(itm)
			case "file":
				resposta, err = b.Process.ProcessaFile(itm)
			default:
				// Se o tipo não for reconhecido, simplesmente retorna
				return
			}

			if err != nil {
				// Tenta enviar o erro, mas não bloqueia se já houver um erro
				select {
				case errCh <- err:
					// Cancela o contexto para outras goroutines

				default:
				}
				return
			}

			// Armazena o resultado no índice correto
			results[idx] = resposta
		}(index, item)
	}

	// Espera pelas goroutines em uma goroutine separada
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// Espera ou por todas as goroutines terminarem ou por um erro ocorrer
	select {
	case <-done:
		// Todas as goroutines completaram sem erros
	case err := <-errCh:
		// Um erro ocorreu, retorna
		return err
	}

	// Une os resultados mantendo a ordem original
	TexUnion := strings.Join(results, " ")

	resp, err := b.IA.SendText(TexUnion + " " + dados.Body[0].JSONString)
	if err != nil {
		return err
	}

	if resp == "" {
		return nil
	}

	TratarTexto, err := utilitariosgorms.ProcessInputText(resp, "url")
	if err != nil {
		return err
	}
	for _, v := range TratarTexto {
		time.Sleep(500 * time.Millisecond)
		err = b.Entrega.EnviaMensagem(v.RespostaIA)
		if err != nil {
			return err
		}
	}

	_ = b.DB.ContagemDeRespostas(dados.Body[0].IDWorkSpace)

	return nil
}

func (b *Backgroud) ProcessaBackgroundFila(dados entities.Dados) error {
	// Preparar slice para armazenar os resultados na ordem correta
	results := make([]string, len(dados.Body)-1)

	// Itera sequencialmente sobre os itens, começando do índice 1
	for i := 1; i < len(dados.Body); i++ {
		index := i - 1 // Índice no slice de resultados
		item := dados.Body[i]

		var resposta string
		var err error

		// Processa o item com base no seu tipo
		switch item.Type {
		case "text":
			resposta, err = b.Process.ProcessaTexto(item)
		case "image":
			resposta, err = b.Process.ProcessaImagem(item)
		case "audio":
			resposta, err = b.Process.ProcessaAudio(item)
		case "file":
			resposta, err = b.Process.ProcessaFile(item)
		default:
			// Se o tipo não for reconhecido, pula para o próximo item
			continue
		}

		// Verifica se houve erro no processamento
		if err != nil {
			return err
		}

		// Armazena o resultado no índice correto
		results[index] = resposta
	}

	// Une os resultados mantendo a ordem original
	TexUnion := strings.Join(results, " ")

	// Envia o texto combinado para a IA
	resp, err := b.IA.SendText(TexUnion)
	if err != nil {
		return err
	}

	// Processa a resposta da IA
	TratarTexto, err := utilitariosgorms.ProcessInputText(resp, "url")
	if err != nil {
		return err
	}

	// Envia cada mensagem processada com um intervalo de 500ms
	for _, v := range TratarTexto {
		time.Sleep(500 * time.Millisecond)
		err = b.Entrega.EnviaMensagem(v.RespostaIA)
		if err != nil {
			return err
		}
	}

	_ = b.DB.ContagemDeRespostas(dados.Body[0].IDWorkSpace)

	return nil
}
