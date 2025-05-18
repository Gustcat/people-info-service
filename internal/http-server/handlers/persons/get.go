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
func GetByID(ctx context.Context, getter Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "hadlers.GetByID"

		id, isParse := params.ParseIDParam(w, r)
		if !isParse {
			return
		}

		//TODO: ошибки 400 и 404

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
