package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
	//"time"
	//"unicode"
)


// ==========================================
// 7. LOOP PRINCIPAL
// ==========================================

func main() {
	fontes := []GameSource{
		MockFonte1{Numero: 1, Arquivo: "mock1.json"},
		MockFonte2{Numero: 2, Arquivo: "mock2.json"},
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\n==========================================")
		fmt.Print("Digite o nome do jogo: ")
		if !scanner.Scan() {
			break
		}
		nomeBusca := strings.TrimSpace(scanner.Text())
		if nomeBusca == "" {
			continue
		}

		fmt.Print("Qual a plataforma (STEAM/PS4/PS5/XBOX ONE/XBOX SERIES S/SWITCH)? ")
		if !scanner.Scan() {
			break
		}
		plataformaBusca := strings.TrimSpace(scanner.Text())

		resultados := executarBuscaConcorrente(fontes, nomeBusca, plataformaBusca)

		if len(resultados) == 0 {
			fmt.Println("\nNão foi possível encontrar nenhum jogo com essas especificações.")
		} else {
			// Determina a oferta de menor valor
			melhor := resultados[0]
			minPreco := math.MaxFloat64

			for _, item := range resultados {
				if item.Preco < minPreco {
					minPreco = item.Preco
					melhor = item
				}
			}

			fmt.Println()
			fmt.Printf("Nome: %s\n", melhor.Nome)
			fmt.Printf("Plataforma: %s\n", melhor.Plataforma)
			fmt.Printf("Valor mais baixo: R$%.2f\n", melhor.Preco)
			fmt.Printf("Encontrado em: Mock Número %d\n", melhor.MockID)
		}

		fmt.Print("\nDeseja buscar outro jogo? (s/n): ")
		if !scanner.Scan() {
			break
		}
		resp := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if resp != "s" && resp != "sim" {
			fmt.Println("Encerrando programa. Até logo!")
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler entrada: %v\n", err)
	}
}