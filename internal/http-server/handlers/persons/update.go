package persons

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"

	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/validation"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-playground/validator/v10"
)

type Updater interface {
	Update(ctx context.Context, id int64, personUpdate *models.PersonUpdate) (*models.FullPerson, error)
}

// Update редактирует профиль человека по ID
//
// @Summary      Редактирует профиль человека
// @Description  У записи с определенным ID редактирует поля
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        id  path      int  true  "Идентификатор профиля человека"
// @Param        input body models.PersonUpdate true "Редактируемые поля"
// @Success      200  {object}  swagger.FullPersonResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      404  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/{id} [patch]
func Update(log *slog.Logger, updater Updater) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.Update"
		log := log.With(slog.String("op", op))

		id, isParse := params.ParseIDParam(c, log)
		if !isParse {
			return
		}

		var personUpdate *models.PersonUpdate

		log.Debug("Receive update request")
		if err := c.ShouldBindJSON(&personUpdate); err != nil {
			if validateErrs, ok := err.(validator.ValidationErrors); ok {
				errMsg := validation.ErrorMessage(validateErrs)
				log.Error("Validation failure", slog.String("error", errMsg))
				c.AbortWithStatusJSON(http.StatusBadRequest, response.Error(errMsg))
				return
			}
			log.Error("Failed to parse request", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("failed to parse request"))
			return
		}

		if isEmptyPersonUpdate(personUpdate) {
			log.Error("Empty request body")
			c.AbortWithStatusJSON(http.StatusBadRequest, response.Error("empty request"))
			return
		}

		log.Debug("Parsed update successfully", slog.Any("person", personUpdate))

		log.Debug("Try to update person in DB")
		person, err := updater.Update(c.Request.Context(), id, personUpdate)
		if err != nil {
			log.Error("Failed to update person", slog.String("error", err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.Error("failed to change person"))
			return
		}

		log.Info("Person updated", slog.Int64("id", id))
		c.JSON(http.StatusOK, response.OK[models.FullPerson](person))
	}
}

func isEmptyPersonUpdate(p *models.PersonUpdate) bool {
	return p.Name == nil &&
		p.Surname == nil &&
		p.Patronymic == nil &&
		p.Age == nil &&
		p.Gender == nil &&
		p.Nationality == nil
}
