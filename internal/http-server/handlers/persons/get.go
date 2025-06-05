package persons

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/Gustcat/people-info-service/internal/repository"
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
func GetByID(log *slog.Logger, getter Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "hadlers.GetByID"
		log := log.With(slog.String("op", op))

		id, isParse := params.ParseIDParam(c, log)
		if !isParse {
			return
		}

		log.Debug("Try to get person", slog.Int64("id", id))
		person, err := getter.GetByID(c.Request.Context(), id)
		if errors.Is(err, repository.ErrPersonNotFound) {
			log.Error("Failed to get person", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusNotFound, response.Error(fmt.Sprintf("Person with id=%d not found", id)))
			return
		}

		if err != nil {
			log.Error("Error calling GetByID", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("failed to get person"))
			return
		}

		c.JSON(http.StatusOK, response.OK[models.FullPerson](person))
	}
}
