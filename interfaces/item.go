package interfaces

// ==========================================
// 1. MODELO CANÔNICO
// ==========================================

type Item struct {
	Nome       string
	Plataforma string
	Preco      float64
	LojaID     int
}

// ==========================================
// 2. INTERFACE
// ==========================================

type GameSource interface {
	ID() int
	Buscar(termo string, plataforma string) ([]Item, error)
}
