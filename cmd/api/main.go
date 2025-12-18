package main

import (
	"fmt"
	"net/http"

	// Ajuste os imports para o SEU caminho do github
	"github.com/GabrielHKGodinho/url-shortener/internal/api"
	"github.com/GabrielHKGodinho/url-shortener/internal/store"
)

func main() {
	connStr := "user=postgres password=secret dbname=shortener sslmode=disable"

	postgresStore, err := store.NewPostgresStore(connStr)
	if err != nil {
		panic(err) // Se não conectar no banco, a API nem deve subir
	}

	// Como PostgresStore implementa a interface Store, o api.Handle aceita ele!
	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", api.HandleShorten(postgresStore))
	mux.HandleFunc("/", api.HandleRedirect(postgresStore))

	fmt.Println("Servidor rodando na porta :8080 (Conectado ao Postgres)")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
