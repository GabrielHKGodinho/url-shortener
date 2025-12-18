package api

import (
	"github.com/GabrielHKGodinho/url-shortener/internal/store"
	"github.com/GabrielHKGodinho/url-shortener/internal/utils"
)

// Service é o nosso "Chef de Cozinha". Ele recebe os pedidos e usa o banco de dados.
type Service struct {
	store store.Store
}

// NewService cria um novo serviço (É essa função que o teste estava procurando e não achava!)
func NewService(s store.Store) *Service {
	return &Service{store: s}
}

// Shorten gera o código curto e salva no banco
func (s *Service) Shorten(originalURL string) (string, error) {
	// 1. Gera o código aleatório
	code := utils.GenerateShortCode(6)

	// 2. Cria o objeto Link
	link := store.Link{
		OriginalURL: originalURL,
		ShortCode:   code,
	}

	// 3. Manda o banco salvar
	if err := s.store.Save(link); err != nil {
		return "", err
	}

	return code, nil
}

// GetOriginal busca o link original e conta o clique
func (s *Service) GetOriginal(code string) (string, error) {
	// 1. Busca no banco
	link, err := s.store.Get(code)
	if err != nil {
		return "", err
	}

	// 2. Incrementa o clique (aquela função nova que criamos)
	_ = s.store.IncrementClick(code)

	return link.OriginalURL, nil
}
