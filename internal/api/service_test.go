package api

import (
	"testing"

	"github.com/GabrielHKGodinho/url-shortener/internal/store"
)

// 1. O MockStore é o nosso "Banco de Mentira".
// Ele guarda os dados num Map simples apenas durante o teste.
type MockStore struct {
	links map[string]store.Link
}

// 2. Implementamos o método Save da interface Store.
// Em vez de ir no Postgres, ele salva no map do Mock.
func (m *MockStore) Save(l store.Link) error {
	m.links[l.ShortCode] = l
	return nil
}

// 3. Implementamos o método Get.
func (m *MockStore) Get(code string) (store.Link, error) {
	// Se o link existir no map, retorna ele.
	if link, ok := m.links[code]; ok {
		return link, nil
	}
	// Se não, retorna vazio (simulando erro de banco)
	return store.Link{}, nil
}

// 4. Implementamos o IncrementClick.
// Como sua Interface agora exige isso, o Mock TEM que ter esse método.
func (m *MockStore) IncrementClick(code string) error {
	// Por enquanto não precisamos fazer nada aqui para o teste de criar link
	return nil
}

func TestShorten(t *testing.T) {
	// --- ARRANGE ---
	// Inicializamos o map (senão dá panic)
	mock := &MockStore{links: make(map[string]store.Link)}

	// Criamos o Service injetando o Mock (em vez do Postgres)
	svc := NewService(mock)

	// --- ACT ---
	originalURL := "https://www.google.com"
	code, err := svc.Shorten(originalURL)

	// --- ASSERT ---
	// 1. Não deve retornar erro
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}

	// 2. O código gerado deve ter 6 caracteres (regra que definimos no utils)
	if len(code) != 6 {
		t.Errorf("Esperava código com 6 caracteres, recebeu: %d", len(code))
	}

	// 3. (Opcional) Verificar se realmente salvou no Mock
	if _, exists := mock.links[code]; !exists {
		t.Error("O link não foi salvo no banco de dados (Mock)")
	}
}

func TestGetOriginal(t *testing.T) {
	// 1. SETUP: Criamos o mock e o service
	mock := &MockStore{links: make(map[string]store.Link)}
	svc := NewService(mock)

	// 2. ARRANGE (Preparar o terreno):
	// Vamos colocar um link "de mentira" dentro do map do Mock manualmente.
	// Assim, quando o service procurar por "codigo123", ele vai achar.
	expectedURL := "https://www.google.com"
	code := "codigo123"

	mock.links[code] = store.Link{
		OriginalURL: expectedURL,
		ShortCode:   code,
	}

	// 3. ACT (Ação): Pedimos para o Service buscar esse código
	resultURL, err := svc.GetOriginal(code)

	// 4. ASSERT (Verificação):
	// Não pode dar erro
	if err != nil {
		t.Errorf("Erro inesperado: %v", err)
	}

	// A URL que voltou tem que ser a mesma que colocamos lá em cima
	if resultURL != expectedURL {
		t.Errorf("Esperava '%s', mas recebeu '%s'", expectedURL, resultURL)
	}
}
