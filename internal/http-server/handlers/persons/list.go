package persons

import (
	"context"
	"github.com/Gustcat/people-info-service/internal/lib/filter"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"log/slog"
	"net/http"
)

type Lister interface {
	List(ctx context.Context, filter *filter.PersonFilter) ([]*models.FullPerson, error)
}

// List возвращает профили людей
//
// @Summary      Возвращает профили людей
// @Description  Возращает профили всех людей с возможностью фильтрации по значению полей и пагинации
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        filter query filter.PersonFilter  false "Фильтрация и пагинация"
// @Success      200  {object}  swagger.FullPersonsResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/ [get]
func List(ctx context.Context, log *slog.Logger, lister Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.List"

		var decoder = schema.NewDecoder()

		log.Debug("Receive list request")
		var personFilter filter.PersonFilter
		err := decoder.Decode(&personFilter, r.URL.Query())
		if err != nil {
			log.Error("Bad request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid query-parameter"))
			return
		}

		log.Debug("Get persons from DB by filter", slog.Any("filter", personFilter))
		persons, err := lister.List(ctx, &personFilter)
		if err != nil {
			log.Error("Failed to list persons", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to list persons"))
			return
		}

		render.JSON(w, r, response.OK[[]*models.FullPerson](&persons))
	}
}
