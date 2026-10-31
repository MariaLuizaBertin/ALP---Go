package fontes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"example.com/m/v2/busca"
	"example.com/m/v2/interfaces"
)

// ==========================================
// IMPLEMENTAÇÃO DO MOCK 6 (mock6.json)
// ==========================================

type MockFonte6 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja6 struct {
	ProductName string  `json:"product_name"`
	Amount      float64 `json:"amount"`
	Hardware    string  `json:"hardware"`
}

func (m MockFonte6) ID() int { return m.Numero }

func (m MockFonte6) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja6
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.Hardware, plataforma) {
			if busca.AtendeBusca(termo, p.ProductName) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.ProductName,
					Plataforma: strings.ToUpper(p.Hardware),
					Preco:      p.Amount,
					MockID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}

