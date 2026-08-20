package handler

import (
	"log"
	"net/http"

	"github.com/12sub/Reset/internal/service"
	"github.com/12sub/Reset/internal/templates"
	"github.com/google/uuid"
)

type TrackedHandler struct {
	service *service.TrackedService
}

func NewTrackedHandler(svc *service.TrackedService) *TrackedHandler {
	return &TrackedHandler{service: svc}
}

func (h *TrackedHandler) TestPage(w http.ResponseWriter, r *http.Request) {
    templates.T.ExecuteTemplate(w, "layout", map[string]string{"Page": "test"})
}

func (h *TrackedHandler) DetectPage(w http.ResponseWriter, r *http.Request) {
    if err := templates.T.ExecuteTemplate(w, "layout", map[string]string{"Page": "detect"}); err != nil {
        log.Printf("template error: %v", err)
    }
}

func (h *TrackedHandler) Detect(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		templates.T.ExecuteTemplate(w, "detect_result", map[string]interface{}{"Error": "Invalid form"})
		return
	}

	email := r.FormValue("email")
	if email == "" {
		email = "demo@reset.ng"
	}
	text := r.FormValue("text")

	sub, err := h.service.Detect(r.Context(), email, text)
	if err != nil {
		templates.T.ExecuteTemplate(w, "detect_result", map[string]interface{}{"Error": err.Error()})
		return
	}

	templates.T.ExecuteTemplate(w, "detect_result", map[string]interface{}{
		"ID":         sub.ID.String(),
		"Merchant":   sub.Merchant,
		"Category":   sub.Category,
		"Amount":     sub.Amount,
		"RawText":    sub.RawText,
		"Status":     sub.Status,
		"Confidence": int(sub.Confidence * 100),
		"MatchType":  sub.MatchType,
	})
}

func (h *TrackedHandler) List(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		email = "demo@reset.ng"
	}

	subs, err := h.service.ListByEmail(r.Context(), email)
	if err != nil {
		renderError(w, err.Error())
		return
	}

	data := map[string]interface{}{
		"Page": "tracked",
		"Subs": subs,
	}

	if err := templates.T.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("template error: %v", err)
	}
}
func (h *TrackedHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		renderError(w, "Invalid form")
		return
	}

	idStr := r.FormValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		renderError(w, "Invalid ID")
		return
	}

	sub, err := h.service.Cancel(r.Context(), id)
	if err != nil {
		renderError(w, err.Error())
		return
	}

	templates.T.ExecuteTemplate(w, "tracked_cancelled", map[string]interface{}{
		"ID":       sub.ID.String(),
		"Merchant": sub.Merchant,
		"Category": sub.Category,
		"Amount":   sub.Amount,
		"Canceled": sub.CanceledAt.Format("Jan 02, 15:04"),
	})
}