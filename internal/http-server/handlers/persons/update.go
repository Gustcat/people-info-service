package persons

import (
	"context"
	"errors"
	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"io"
	"log"
	"net/http"
)

type Updater interface {
	Update(ctx context.Context, id int64, personUpdate *models.PersonUpdate) (*models.FullPerson, error)
}

func Update(ctx context.Context, updater Updater) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Update"

		id, isParse := params.ParseIDParam(w, r)
		if !isParse {
			return
		}

		var personUpdate *models.PersonUpdate

		err := render.DecodeJSON(r.Body, &personUpdate)
		if errors.Is(err, io.EOF) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to parse request"))
			return
		}

		person, err := updater.Update(ctx, id, personUpdate)
		if err != nil {
			log.Printf("failed to update person: %s", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to change person"))
			return
		}

		render.JSON(w, r, response.OK[models.FullPerson](person))
	})
}
