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
// IMPLEMENTAÇÃO DO MOCK 3 (loja3.json)
// ==========================================

type LojaFonte3 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja3 struct {
	Nome       string  `json:"nome"`
	Preco      float64 `json:"preco"`
	Plataforma string  `json:"plataforma"`
}

func (m LojaFonte3) ID() int { return m.Numero }

func (m LojaFonte3) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja3
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.Plataforma, plataforma) || strings.EqualFold(plataforma, "TODAS") {
			if busca.AtendeBusca(termo, p.Nome) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.Nome,
					Plataforma: strings.ToUpper(p.Plataforma),
					Preco:      p.Preco,
					LojaID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}
