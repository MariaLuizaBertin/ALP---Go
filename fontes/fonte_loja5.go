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
// IMPLEMENTAÇÃO DO MOCK 5 (loja5.json)
// ==========================================

type LojaFonte5 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja5 struct {
	Jogo             string  `json:"jogo"`
	PrecoPromocional float64 `json:"preco_promocional"`
	Dispositivo      string  `json:"dispositivo"`
}

func (m LojaFonte5) ID() int { return m.Numero }

func (m LojaFonte5) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja5
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.Dispositivo, plataforma) || strings.EqualFold(plataforma, "TODAS") {
			if busca.AtendeBusca(termo, p.Jogo) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.Jogo,
					Plataforma: strings.ToUpper(p.Dispositivo),
					Preco:      p.PrecoPromocional,
					LojaID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}
