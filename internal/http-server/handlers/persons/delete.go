package persons

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
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
func Delete(log *slog.Logger, deleter Deleter) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.Delete"
		log := log.With(slog.String("op", op))

		id, isParse := params.ParseIDParam(c, log)
		if !isParse {
			return
		}

		if err := deleter.Delete(c.Request.Context(), id); err != nil {
			log.Error("Failed to delete person", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("failed to delete person"))
			return
		}

		log.Info("Person deleted", slog.Int64("id", id))
		// можно со статусом 204 обработать вариант, когда совершается попытка удалить несуществующий объект
		c.JSON(http.StatusOK, response.OK[struct{}](nil))
	}
}
