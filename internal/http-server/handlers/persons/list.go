package persons

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/filter"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

type Lister interface {
	List(ctx context.Context, filter *filter.PersonFilter) ([]*models.FullPerson, uint64, error)
}

// List возвращает профили людей
//
// @Summary      Возвращает профили людей
// @Description  Возращает профили всех людей с возможностью фильтрации по значению полей и пагинации
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        filter query filter.PersonFilter  false "Фильтрация и пагинация"
// @Success      200  {object}  swagger.PersonsWithPaginationResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/ [get]
func List(ctx context.Context, log *slog.Logger, lister Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.List"

		decoder := schema.NewDecoder()

		log.Debug("Receive list request")
		var personFilter filter.PersonFilter
		err := decoder.Decode(&personFilter, r.URL.Query())
		if err != nil {
			log.Error("Bad request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(fmt.Sprintf("invalid query-parameter: %s", err.Error())))
			return
		}

		log.Debug("Get persons from DB by filter", slog.Any("filter", personFilter))
		persons, total, err := lister.List(ctx, &personFilter)
		if err != nil {
			log.Error("Failed to list persons", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to list persons"))
			return
		}

		url := urlbuilder.BaseURL(r)
		offset := response.DefaultOffset
		limit := response.DefaultLimit
		if personFilter.Limit != nil {
			limit = *personFilter.Limit
		}
		if personFilter.Offset != nil {
			offset = *personFilter.Offset
		}
		pagination, err := response.NewPagination(limit, offset, total, url)
		if err != nil {
			log.Error("Failed to create pagination", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to create pagination"))
		}

		render.JSON(w, r, response.OKWithPagination[[]*models.FullPerson](&persons, pagination))
	}
}
