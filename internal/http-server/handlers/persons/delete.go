package persons

import (
	"context"
	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/go-chi/render"
	"net/http"
)

type Deleter interface {
	Delete(ctx context.Context, id int64) error
}

func Delete(ctx context.Context, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Delete"

		id, isParse := params.ParseIDParam(w, r)
		if !isParse {
			return
		}

		if err := deleter.Delete(r.Context(), id); err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete person"))
			return
		}

		// можно со статусом 204 обработать вариант, когда совершается попытка удалить несуществующий объект
		render.JSON(w, r, response.OK[struct{}](nil))
	}
}
