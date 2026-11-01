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
// 4. IMPLEMENTAÇÃO DO MOCK 2 (loja2.json)
// ==========================================

type LojaFonte2 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja2 struct {
	GameName string  `json:"game_name"`
	Price    float64 `json:"price"`
	System   string  `json:"system"`
}

func (m LojaFonte2) ID() int { return m.Numero }

func (m LojaFonte2) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja2
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.System, plataforma) || strings.EqualFold(plataforma, "TODAS") {
			if busca.AtendeBusca(termo, p.GameName) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.GameName,
					Plataforma: strings.ToUpper(p.System),
					Preco:      p.Price,
					LojaID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}
