package persons

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"github.com/Gustcat/people-info-service/internal/lib/validation"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/Gustcat/people-info-service/internal/repository"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const (
	urlForAge         = "https://api.agify.io/"
	urlForGender      = "https://api.genderize.io/"
	urlForNationality = "https://api.nationalize.io/"
	nameParam         = "name"
)

type Creator interface {
	Create(ctx context.Context, person *models.EnrichmentPerson) (int64, error)
}

// Create создает профиль человека
//
// @Summary      Создать профиль человека
// @Description  Вводится ФИО, данные обогащаются возрастом, национальностью и полом, возращается ID созданной записи
// @Tags         persons
// @Accept       json
// @Produce      json
// @Param        input body models.Person true "ФИО"
// @Success      201  {object}  swagger.IdResponse
// @Failure      400  {object}  swagger.ErrorResponse
// @Failure      500  {object}  swagger.ErrorResponse
// @Router       /persons/ [post]
func Create(ctx context.Context, log *slog.Logger, creator Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Create"
		log := log.With(slog.String("op", op))

		var person models.Person

		log.Debug("Receive create request")
		err := render.DecodeJSON(r.Body, &person)
		if errors.Is(err, io.EOF) {
			log.Error("Bad request: empty request")
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("empty request"))
			return
		}

		if err != nil {
			log.Error("Failed to parse request", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error("failed to parse request"))
			return
		}
		log.Debug("Parsed create successfully", slog.Any("person", person))

		if err := validator.New().Struct(person); err != nil {
			validateErr := err.(validator.ValidationErrors)
			errMsg := validation.ErrorMessage(validateErr)
			log.Error("Validation failure", slog.String("error", errMsg))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(errMsg))
			return
		}

		log.Debug("Try to enrich person information")
		enrichPerson := enrichPerson(&person, log)
		log.Debug("Enrich person successfully", slog.Any("enrich", enrichPerson))

		id, err := creator.Create(r.Context(), enrichPerson)
		if errors.Is(err, repository.ErrPersonExists) {
			log.Error("Get error", slog.String("error", err.Error()))
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(fmt.Sprintf(
				"Person with name %s %s already exists", person.Name, person.Surname)))
			return
		}

		if err != nil {
			log.Error("Failed to add person", slog.String("error", err.Error()))
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add person"))
			return
		}

		log.Info("Person created", slog.Int64("id", id))
		createResp := &models.Identifier{ID: id}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.OK[models.Identifier](createResp))
	}
}

func enrichPerson(person *models.Person, log *slog.Logger) *models.EnrichmentPerson {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

	enrichPerson := &models.EnrichmentPerson{
		Person: *person,
	}

	wg.Add(3)

	go func() {
		defer wg.Done()
		fulUrlForAge, err := urlbuilder.BuildWithQueryParams(urlForAge, map[string]string{nameParam: person.Name})
		if err != nil {
			log.Error("Failed to parse url for age", slog.String("error", err.Error()))
			return
		}
		log.Debug("URL for age", slog.String("url", fulUrlForAge))

		resp, err := http.Get(fulUrlForAge)
		if err != nil {
			log.Error("Failed to fetch url for age", slog.String("error", err.Error()))
			return
		}
		defer resp.Body.Close()

		var data struct {
			Age *int64 `json:"age"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Error("Failed to decode response body", slog.String("error", err.Error()))
			return
		}
		if data.Age == nil {
			log.Debug("Receive age", slog.Int64("age", *data.Age))
		} else {
			log.Debug("Receive empty age")
		}

		mu.Lock()
		enrichPerson.Age = data.Age
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		fulUrlForGender, err := urlbuilder.BuildWithQueryParams(urlForGender, map[string]string{nameParam: person.Name})
		if err != nil {
			log.Error("Failed to parse url for gender", slog.String("error", err.Error()))
		}
		log.Debug("URL for gender", slog.String("url", fulUrlForGender))

		resp, err := http.Get(fulUrlForGender)
		if err != nil {
			log.Error("Failed to fetch url for gender", slog.String("error", err.Error()))
		}
		defer resp.Body.Close()

		var data struct {
			Gender      *models.Gender `json:"gender"`
			Probability float64        `json:"probability"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Error("Failed to decode response body", slog.String("error", err.Error()))
		}
		log.Debug("Receive answer", slog.Any("data", data))
		if data.Probability < 0.7 {
			data.Gender = nil
		}

		mu.Lock()
		enrichPerson.Gender = data.Gender
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		fullUrlNationality, err := urlbuilder.BuildWithQueryParams(urlForNationality, map[string]string{nameParam: person.Name})
		if err != nil {
			log.Error("Failed to parse url for nationality", slog.String("error", err.Error()))
		}
		log.Debug("URL for nationality", slog.String("url", fullUrlNationality))

		resp, err := http.Get(fullUrlNationality)
		if err != nil {
			log.Error("Failed to fetch url for nationality", slog.String("error", err.Error()))
		}
		defer resp.Body.Close()

		var data struct {
			Country []struct {
				CountryID   string  `json:"country_id"`
				Probability float64 `json:"probability"`
			} `json:"country"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Error("Failed to decode response body", slog.String("error", err.Error()))
		}
		log.Debug("Get countries", slog.Int("count", len(data.Country)))

		var nationality *string
		probability := 0.0

		for _, version := range data.Country {
			if version.Probability > probability {
				nationality = &version.CountryID
				probability = version.Probability
			}
		}

		mu.Lock()
		enrichPerson.Nationality = nationality
		mu.Unlock()
	}()

	wg.Wait()

	return enrichPerson
}
