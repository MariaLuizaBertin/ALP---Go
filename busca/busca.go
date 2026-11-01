package busca

import (
	"errors"
	"fmt"
	"io/fs"
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

// MensagemAvisoLoja gera uma mensagem de aviso amigável para o usuário com base no erro retornado pela loja.
func MensagemAvisoLoja(id int, err error) string {
	if err == nil {
		return ""
	}

	motivo := strings.ToLower(err.Error())

	switch {
	case errors.Is(err, fs.ErrNotExist) ||
		strings.Contains(motivo, "no such file") ||
		strings.Contains(motivo, "not exist") ||
		strings.Contains(motivo, "cannot find the file specified") ||
		strings.Contains(motivo, "arquivo não encontrado"):
		return fmt.Sprintf("A loja %d não pode ser consultada porque o arquivo não foi encontrado.", id)
	case errors.Is(err, fs.ErrPermission) ||
		strings.Contains(motivo, "permission denied") ||
		strings.Contains(motivo, "acesso negado") ||
		strings.Contains(motivo, "não pode ser aberto") ||
		strings.Contains(motivo, "forbidden"):
		return fmt.Sprintf("A loja %d não pode ser consultada porque o arquivo não pode ser aberto.", id)
	case strings.Contains(motivo, "connection refused") || strings.Contains(motivo, "timeout") || strings.Contains(motivo, "dial tcp") || strings.Contains(motivo, "network is unreachable") || strings.Contains(motivo, "connection reset"):
		return fmt.Sprintf("A loja %d não pode ser consultada devido a uma falha na conexão.", id)
	default:
		return fmt.Sprintf("A loja %d não pode ser consultada no momento. Motivo: %s", id, err)
	}
}

// ==========================================
// 6. BUSCA CONCORRENTE VIA GOROUTINES
// ==========================================

func ExecutarBuscaConcorrente(fontes []interfaces.GameSource, termo, plataforma string) ([]interfaces.Item, []string) {
	canal := make(chan []interfaces.Item, len(fontes))
	avisos := make(chan string, len(fontes))
	var wg sync.WaitGroup

	for _, fonte := range fontes {
		wg.Add(1)
		go func(src interfaces.GameSource) {
			defer wg.Done()
			itens, err := src.Buscar(termo, plataforma)
			if err != nil {
				avisos <- MensagemAvisoLoja(src.ID(), err)
				return
			}
			canal <- itens
		}(fonte)
	}

	go func() {
		wg.Wait()
		close(canal)
		close(avisos)
	}()

	var consolidados []interfaces.Item
	for lista := range canal {
		consolidados = append(consolidados, lista...)
	}

	var avisosLista []string
	for aviso := range avisos {
		avisosLista = append(avisosLista, aviso)
	}

	return consolidados, avisosLista
}
