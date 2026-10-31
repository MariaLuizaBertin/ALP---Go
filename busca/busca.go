package busca

import (
	"strings"
	"sync"

	"example.com/m/v2/interfaces"
)

// ==========================================
// 5. ALGORITMO COM RUNAS (SUPORTE UTF-8 / ACENTOS)
// ==========================================

func AtendeBusca(termo, candidato string) bool {
	tNorm := strings.ToLower(termo)
	cNorm := strings.ToLower(candidato)

	if strings.Contains(cNorm, tNorm) {
		return true
	}
	return CalcularSimilaridade(tNorm, cNorm) > 0.35
}

func CalcularSimilaridade(s1, s2 string) float64 {
	// Converte strings para slices de runas para não quebrar caracteres acentuados
	r1 := []rune(s1)
	r2 := []rune(s2)

	getTrigrams := func(runas []rune) map[string]bool {
		trigrams := make(map[string]bool)
		if len(runas) < 3 {
			// Se o texto tiver menos de 3 runas, usa a palavra inteira
			if len(runas) > 0 {
				trigrams[string(runas)] = true
			}
			return trigrams
		}
		for i := 0; i <= len(runas)-3; i++ {
			// Fatia as runas com segurança e gera a chave
			trigrams[string(runas[i:i+3])] = true
		}
		return trigrams
	}

	t1 := getTrigrams(r1)
	t2 := getTrigrams(r2)

	if len(t1) == 0 || len(t2) == 0 {
		return 0
	}

	interseccao := 0
	for k := range t1 {
		if t2[k] {
			interseccao++
		}
	}

	return 2.0 * float64(interseccao) / float64(len(t1)+len(t2))
}

// ==========================================
// 6. BUSCA CONCORRENTE VIA GOROUTINES
// ==========================================

func ExecutarBuscaConcorrente(fontes []interfaces.GameSource, termo, plataforma string) []interfaces.Item {
	canal := make(chan []interfaces.Item, len(fontes))
	var wg sync.WaitGroup

	for _, fonte := range fontes {
		wg.Add(1)
		go func(src interfaces.GameSource) {
			defer wg.Done()
			itens, err := src.Buscar(termo, plataforma)
			if err == nil {
				canal <- itens
			}
		}(fonte)
	}

	go func() {
		wg.Wait()
		close(canal)
	}()

	var consolidados []interfaces.Item
	for lista := range canal {
		consolidados = append(consolidados, lista...)
	}
	return consolidados
}

