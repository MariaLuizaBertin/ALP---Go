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
// IMPLEMENTAÇÃO DO MOCK 4 (loja4.json)
// ==========================================

type LojaFonte4 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja4 struct {
	Title    string  `json:"title"`
	Cost     float64 `json:"cost"`
	Platform string  `json:"platform"`
}

func (m LojaFonte4) ID() int { return m.Numero }

func (m LojaFonte4) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja4
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.Platform, plataforma) || strings.EqualFold(plataforma, "TODAS") {
			if busca.AtendeBusca(termo, p.Title) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.Title,
					Plataforma: strings.ToUpper(p.Platform),
					Preco:      p.Cost,
					LojaID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}
