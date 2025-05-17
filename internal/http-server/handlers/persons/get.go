package persons

import (
	"context"
	"github.com/Gustcat/people-info-service/internal/lib/params"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"log"
	"net/http"
)

type Getter interface {
	GetByID(ctx context.Context, id int64) (*models.FullPerson, error)
}

func GetByID(ctx context.Context, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hadlers.GetByID"

		id, isParse := params.ParseIDParam(w, r)
		if !isParse {
			return
		}

		person, err := getter.GetByID(ctx, id)
		if err != nil {
			log.Printf("error calling GetByID: %s", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to get person"))
			return
		}

		render.JSON(w, r, response.OK[models.FullPerson](person))
	}
}
