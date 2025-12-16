package main

import (
	"fmt"
	"net/http"

	// Ajuste os imports para o SEU caminho do github
	"github.com/GabrielHKGodinho/url-shortener/internal/api"
	"github.com/GabrielHKGodinho/url-shortener/internal/store"
)

func main() {
	// 1. Inicializa o Banco
	memStore := store.NewMemoryStore()

	// 2. Configura Rotas
	// Perceba como passamos o 'memStore' para dentro dos handlers.
	// Isso se chama INJEÇÃO DE DEPENDÊNCIA (Manual).
	mux := http.NewServeMux()

	mux.HandleFunc("/shorten", api.HandleShorten(memStore))
	mux.HandleFunc("/", api.HandleRedirect(memStore))

	// 3. Start
	fmt.Println("Servidor rodando na porta :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
