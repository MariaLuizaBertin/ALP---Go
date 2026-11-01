package main

import (
	"math"
	"strings"

	"example.com/m/v2/busca"
	"example.com/m/v2/fontes"
	"example.com/m/v2/interfaces"
	ui "github.com/webui-dev/go-webui/v2"
)

// ==========================================
// ESTRUTURAS DE RESPOSTA PARA A INTERFACE
// ==========================================

type RespostaBusca struct {
	Encontrado bool              `json:"encontrado"`
	Mensagem   string            `json:"mensagem,omitempty"`
	Avisos     []string          `json:"avisos,omitempty"`
	Melhor     *interfaces.Item  `json:"melhor,omitempty"`
	Resultados []interfaces.Item `json:"resultados,omitempty"`
}

// ==========================================
// LOOP PRINCIPAL COM WEBUI
// ==========================================

func main() {
	// Registro das classes concretas na pasta 'fontes' usando caminhos relativos diretos
	fontes := []interfaces.GameSource{
		fontes.LojaFonte1{Numero: 1, Arquivo: "lojas-de-jogos/loja1.json"},
		fontes.LojaFonte2{Numero: 2, Arquivo: "lojas-de-jogos/loja2.json"},
		fontes.LojaFonte3{Numero: 3, Arquivo: "lojas-de-jogos/loja3.json"},
		fontes.LojaFonte4{Numero: 4, Arquivo: "lojas-de-jogos/loja4.json"},
		fontes.LojaFonte5{Numero: 5, Arquivo: "lojas-de-jogos/loja5.json"},
		fontes.LojaFonte6{Numero: 6, Arquivo: "lojas-de-jogos/loja6.json"},
		fontes.LojaFonte7{Numero: 7, Arquivo: "lojas-de-jogos/loja7.json"},
	}

	// 1. Cria uma nova janela do WebUI
	janela := ui.NewWindow()

	// 2. Ajusta as dimensões e centraliza na tela do PC
	janela.SetSize(1280, 840)
	janela.SetCenter()

	// 3. Conecta a função de busca Go ao JavaScript do frontend
	ui.Bind(janela, "buscarJogo", func(e ui.Event) RespostaBusca {
		nomeBusca, _ := ui.GetArgAt[string](e, 0)
		plataformaBusca, _ := ui.GetArgAt[string](e, 1)

		nomeBusca = strings.TrimSpace(nomeBusca)
		plataformaBusca = strings.TrimSpace(plataformaBusca)

		if nomeBusca == "" || plataformaBusca == "" {
			return RespostaBusca{
				Encontrado: false,
				Mensagem:   "Por favor, informe o nome do jogo e a plataforma.",
			}
		}

		// Executa a busca concorrente em todas as fontes
		resultados, avisos := busca.ExecutarBuscaConcorrente(fontes, nomeBusca, plataformaBusca)

		if len(resultados) == 0 {
			return RespostaBusca{
				Encontrado: false,
				Mensagem:   "Não foi possível encontrar nenhum jogo com essas especificações.",
				Avisos:     avisos,
			}
		}

		// Determina a oferta de menor valor
		melhor := resultados[0]
		minPreco := math.MaxFloat64

		for _, item := range resultados {
			if item.Preco < minPreco {
				minPreco = item.Preco
				melhor = item
			}
		}

		return RespostaBusca{
			Encontrado: true,
			Melhor:     &melhor,
			Resultados: resultados,
			Avisos:     avisos,
		}
	})

	// 4. Configura a pasta da interface web e exibe a página principal
	janela.SetRootFolder("ui")
	janela.Show("index.html")

	// 5. Aguarda até que todas as janelas sejam fechadas
	ui.Wait()
}
