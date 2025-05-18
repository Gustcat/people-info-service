package persons

import (
	"context"
	"errors"
	"fmt"
	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/validation"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log/slog"
	"net/http"
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
func Update(ctx context.Context, log *slog.Logger, updater Updater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Update"

		id, isParse := params.ParseIDParam(w, r, log)
		if !isParse {
			return
		}

		var personUpdate *models.PersonUpdate
		log.Debug("Receive update request")
		err := validation.DecodeStrictJSON(r, &personUpdate)
		if errors.Is(err, io.EOF) || isEmptyPersonUpdate(personUpdate) {
			log.Error("Empty request body")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("Bad request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(fmt.Sprintf("malformed JSON: %s", err)))
			return
		}
		log.Debug("Parsed update successfully", slog.Any("parsed", personUpdate))

		if err := validator.New().Struct(personUpdate); err != nil {
			validateErr := err.(validator.ValidationErrors)
			errMsg := validation.ErrorMessage(validateErr)
			log.Error("Bad request", slog.String("error", errMsg))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(errMsg))
			return
		}

		log.Debug("Try to update person in DB")
		person, err := updater.Update(ctx, id, personUpdate)
		if err != nil {
			log.Error("Failed to update person", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to change person"))
			return
		}

		log.Info("Person updated", slog.Int64("id", id))
		render.JSON(w, r, response.OK[models.FullPerson](person))
	})
}

func isEmptyPersonUpdate(p *models.PersonUpdate) bool {
	return p.Name == nil &&
		p.Surname == nil &&
		p.Patronymic == nil &&
		p.Age == nil &&
		p.Gender == nil &&
		p.Nationality == nil
}
