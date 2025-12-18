package api

import (
	"encoding/json"
	"net/http"

	// Importe seus pacotes
	"github.com/GabrielHKGodinho/url-shortener/internal/store"
	"github.com/GabrielHKGodinho/url-shortener/internal/utils"
)

// Estrutura para os inputs e outputs (DTOs - Data Transfer Objects)
type CreateLinkRequest struct {
	URL string `json:"url"`
}

type CreateLinkResponse struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

// EnviarJSON é um helper para padronizar respostas (DRY - Don't Repeat Yourself)
func sendJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

// EnviarError padroniza erros JSON
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

// HandleShorten cuida da criação do link.
// Note que precisamos passar o 'store' como parâmetro, senão ele não sabe onde salvar.
func HandleShorten(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			sendError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		var req CreateLinkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, http.StatusBadRequest, "Invalid JSON")
			return
		}

		// Regra de Negócio
		code := utils.GenerateShortCode(6)
		link := store.Link{
			OriginalURL: req.URL,
			ShortCode:   code,
		}

		if err := db.Save(link); err != nil {
			sendError(w, http.StatusInternalServerError, "Failed to save link")
			return
		}

		// Resposta
		resp := CreateLinkResponse{
			OriginalURL: link.OriginalURL,
			ShortURL:    "http://localhost:8080/" + code,
		}
		sendJSON(w, http.StatusCreated, resp)
	}
}

// HandleRedirect cuida do redirecionamento.
func HandleRedirect(db store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Pega o código da URL: /abc -> abc
		// Cuidado: Em produção usaríamos um router melhor (Chi ou Gin) para pegar params.
		if len(r.URL.Path) < 2 {
			sendError(w, http.StatusBadRequest, "Invalid code")
			return
		}
		code := r.URL.Path[1:]

		link, err := db.Get(code)
		if err != nil {
			sendError(w, http.StatusNotFound, "Link not found")
			return
		}

		db.IncrementClick(code)

		http.Redirect(w, r, link.OriginalURL, http.StatusTemporaryRedirect)
	}
}
