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
// IMPLEMENTAÇÃO DO MOCK 7 (mock7.json)
// ==========================================

type MockFonte7 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja7 struct {
	ItemTitle  string  `json:"item_title"`
	ValueBRL   float64 `json:"value_brl"`
	TagConsole string  `json:"tag_console"`
}

func (m MockFonte7) ID() int { return m.Numero }

func (m MockFonte7) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja7
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.TagConsole, plataforma) {
			if busca.AtendeBusca(termo, p.ItemTitle) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.ItemTitle,
					Plataforma: strings.ToUpper(p.TagConsole),
					Preco:      p.ValueBRL,
					MockID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}

