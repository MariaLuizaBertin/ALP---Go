package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)


// ==========================================
// 4. IMPLEMENTAÇÃO DO MOCK 2 (mock2.json)
// ==========================================

type MockFonte2 struct {
	Numero  int
	Arquivo string
}

type ProdutoLoja2 struct {
	GameName string  `json:"game_name"`
	Price    float64 `json:"price"`
	System   string  `json:"system"`
}

func (m MockFonte2) ID() int { return m.Numero }

func (m MockFonte2) Buscar(termo string, plataforma string) ([]Item, error) {
	bytes, err := os.ReadFile(m.Arquivo)
	if err != nil {
		return nil, fmt.Errorf("erro lendo %s: %w", m.Arquivo, err)
	}

	var produtos []ProdutoLoja2
	if err := json.Unmarshal(bytes, &produtos); err != nil {
		return nil, fmt.Errorf("erro no parse de %s: %w", m.Arquivo, err)
	}

	var encontrados []Item
	for _, p := range produtos {
		if strings.EqualFold(p.System, plataforma) {
			if atendeBusca(termo, p.GameName) {
				encontrados = append(encontrados, Item{
					Nome:       p.GameName,
					Plataforma: strings.ToUpper(p.System),
					Preco:      p.Price,
					MockID:     m.Numero,
				})
			}
		}
	}
	return encontrados, nil
}
