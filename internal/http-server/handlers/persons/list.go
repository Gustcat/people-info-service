package persons

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/filter"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"github.com/Gustcat/people-info-service/internal/models"
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
func List(log *slog.Logger, lister Lister) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.List"
		log := log.With(slog.String("op", op))

		log.Debug("Receive list request")
		var personFilter filter.PersonFilter
		err := c.BindQuery(&personFilter)
		if err != nil {
			log.Error("Bad request", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(fmt.Sprintf("invalid query-parameter: %s", err.Error())))
			return
		}

		log.Debug("Get persons from DB by filter", slog.Any("filter", personFilter))
		persons, total, err := lister.List(c.Request.Context(), &personFilter)
		if err != nil {
			log.Error("Failed to list persons", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("failed to list persons"))
			return
		}

		url := urlbuilder.BaseURL(c.Request)
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
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("failed to create pagination"))
		}

		c.JSON(http.StatusOK, response.OKWithPagination[[]*models.FullPerson](&persons, pagination))
	}
}
