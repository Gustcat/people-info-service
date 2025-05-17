package persons

import (
	"context"
	"github.com/Gustcat/people-info-service/internal/lib/filter"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
	"log"
	"net/http"
)

type Lister interface {
	List(ctx context.Context, filter *filter.PersonFilter) ([]*models.FullPerson, error)
}

func List(ctx context.Context, lister Lister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.List"

		var decoder = schema.NewDecoder()

		var personFilter filter.PersonFilter
		err := decoder.Decode(&personFilter, r.URL.Query())
		if err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid query-parameter"))
			return
		}

		persons, err := lister.List(ctx, &personFilter)
		if err != nil {
			log.Printf("failed to list persons: %s", err)
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to list persons"))
			return
		}

		render.JSON(w, r, response.OK[[]*models.FullPerson](&persons))
	}
}
