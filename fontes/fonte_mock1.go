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
// 3. IMPLEMENTAÇÃO DO MOCK 1 (mock1.json)
// ==========================================

type MockFonte1 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja1 struct {
	Titulo  string  `json:"titulo"`
	Valor   float64 `json:"valor"`
	Console string  `json:"console"`
}

func (m MockFonte1) ID() int { return m.Numero }

func (m MockFonte1) Buscar(termo string, plataforma string) ([]interfaces.Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja1
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []interfaces.Item
	for _, p := range produtos {
		if strings.EqualFold(p.Console, plataforma) {
			if busca.AtendeBusca(termo, p.Titulo) {
				encontrados = append(encontrados, interfaces.Item{
					Nome:       p.Titulo,
					Plataforma: strings.ToUpper(p.Console),
					Preco:      p.Valor,
					MockID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}

