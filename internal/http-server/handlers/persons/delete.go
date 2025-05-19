package persons

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/go-chi/render"
)

type Deleter interface {
	Delete(ctx context.Context, id int64) error
}

// Delete удаляет профиль человека по ID
//
// @Summary      Удаляет профиль человека
// @Description  Удаляет профиль человека
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Идентификатор профиля человека"
// @Success      200  {object}  swagger.EmptyResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/{id} [delete]
func Delete(ctx context.Context, log *slog.Logger, deleter Deleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Delete"

		id, isParse := params.ParseIDParam(w, r, log)
		if !isParse {
			return
		}

		if err := deleter.Delete(r.Context(), id); err != nil {
			log.Error("Failed to delete person", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to delete person"))
			return
		}

		log.Info("Person deleted", slog.Int64("id", id))
		// можно со статусом 204 обработать вариант, когда совершается попытка удалить несуществующий объект
		render.JSON(w, r, response.OK[struct{}](nil))
	}
}
