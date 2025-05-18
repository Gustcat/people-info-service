package params

import (
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

func ParseIDParam(w http.ResponseWriter, r *http.Request, log *slog.Logger) (int64, bool) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		log.Error("Missing id parameter")
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("missing id parameter"))
		return 0, false
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Invalid id parameter", slog.String("parameter", idStr))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid id parameter"))
		return 0, false
	}

	return id, true
}
