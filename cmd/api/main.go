package main

import (
	"fmt"
	"net/http"

	"github.com/GabrielHKGodinho/url-shortener/internal/store"
)

func main() {
	memStore := store.NewMemoryStore()

	mux := http.NewServeMux()

	// Define uma rota GET "/"
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		link := store.Link{
			ShortCode:   "golang",
			OriginalURL: "https://golang.org",
		}
		memStore.Save(link)

		rec, erro := memStore.Get(link.ShortCode)

		if erro != nil {
			http.Error(w, "Link não encontrado!", 404)
		}

		fmt.Fprintf(w, "Recuperado: %s --> %s", rec.ShortCode, rec.OriginalURL)
	})

	// Configuração do servidor
	port := ":8080"
	fmt.Printf("Servidor rodando na porta %s\n", port)

	// Sobe o servidor. Se der erro, o programa fecha (Panic)
	if err := http.ListenAndServe(port, mux); err != nil {
		panic(err)
	}
}
