package persons

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/Gustcat/people-info-service/internal/repository"
	"github.com/go-chi/render"
)

type Getter interface {
	GetByID(ctx context.Context, id int64) (*models.FullPerson, error)
}

// GetByID возвращает профиль человека по ID
//
// @Summary      Получить профиль человека
// @Description  Возвращает информацию по ID
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Идентификатор профиля человека"
// @Success      200  {object}  swagger.FullPersonResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      404  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/{id} [get]
func GetByID(ctx context.Context, log *slog.Logger, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hadlers.GetByID"
		log := log.With(slog.String("op", op))

		id, isParse := params.ParseIDParam(w, r, log)
		if !isParse {
			return
		}

		log.Debug("Try to get person", slog.Int64("id", id))
		person, err := getter.GetByID(ctx, id)
		if errors.Is(err, repository.ErrPersonNotFound) {
			log.Error("Failed to get person", slog.String("error", err.Error()))
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, response.Error(fmt.Sprintf("Person with id=%d not found", id)))
			return
		}

		if err != nil {
			log.Error("Error calling GetByID", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get person"))
			return
		}

		render.JSON(w, r, response.OK[models.FullPerson](person))
	}
}
