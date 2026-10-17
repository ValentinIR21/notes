// транспортный слой: HTTP обертка вокруг сценариев.

package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"notes/internal/domain"
)

type NoteService interface {
	Create(title, text string) (string, error)
	Get(id string) (domain.Note, error)
	List() ([]domain.Note, error)
	Delete(id string) error
}

type Handler struct {
	svc NoteService
}

func New(svc NoteService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /notes", h.create)

	return recover(Logging(mux))
}

type CreateRequest struct {
	Title string `json: "title"`
	Text  string `json: "text"`
}

type noteResponse struct {
	ID    string `json: "id"`
	Title string `json: "title"`
	Text  string `json: "text"`
}

func toResponse(n domain.Note) noteResponse {
	return noteResponse{ID: n.ID, Title: n.Title, Text: n.Text}
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req CreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	id, err := h.svc.Create(req.Title, req.Text)
	if err != nil {
		writeDomainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"id": id})
}

func writeDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		writeError(w, http.StatusNotFound, "Заметка не найдена")

	case errors.Is(err, domain.ErrEmptyTitle):
		writeError(w, http.StatusBadRequest, "Заголовок обязателен")
	default:
		log.Printf("internal error: %v", err)
		writeError(w, http.StatusInternalServerError, "внутреняя ошибок")
	}
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}
