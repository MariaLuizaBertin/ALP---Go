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
	Melhor     *interfaces.Item  `json:"melhor,omitempty"`
	Resultados []interfaces.Item `json:"resultados,omitempty"`
}

// ==========================================
// LOOP PRINCIPAL COM WEBUI
// ==========================================

func main() {
	// Registro das classes concretas na pasta 'fontes' usando caminhos relativos diretos
	fontes := []interfaces.GameSource{
		fontes.MockFonte1{Numero: 1, Arquivo: "mocks/mock1.json"},
		fontes.MockFonte2{Numero: 2, Arquivo: "mocks/mock2.json"},
		fontes.MockFonte3{Numero: 3, Arquivo: "mocks/mock3.json"},
		fontes.MockFonte4{Numero: 4, Arquivo: "mocks/mock4.json"},
		fontes.MockFonte5{Numero: 5, Arquivo: "mocks/mock5.json"},
		fontes.MockFonte6{Numero: 6, Arquivo: "mocks/mock6.json"},
		fontes.MockFonte7{Numero: 7, Arquivo: "mocks/mock7.json"},
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
		resultados := busca.ExecutarBuscaConcorrente(fontes, nomeBusca, plataformaBusca)

		if len(resultados) == 0 {
			return RespostaBusca{
				Encontrado: false,
				Mensagem:   "Não foi possível encontrar nenhum jogo com essas especificações.",
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
		}
	})

	// 4. Configura a pasta da interface web e exibe a página principal
	janela.SetRootFolder("ui")
	janela.Show("index.html")

	// 5. Aguarda até que todas as janelas sejam fechadas
	ui.Wait()
}
