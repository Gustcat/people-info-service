package persons

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"github.com/Gustcat/people-info-service/internal/lib/validation"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/Gustcat/people-info-service/internal/repository"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
	"sync"
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
func Create(ctx context.Context, creator Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Create"

		var person models.Person

		err := render.DecodeJSON(r.Body, &person)
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

		if err := validator.New().Struct(person); err != nil {
			validateErr := err.(validator.ValidationErrors)
			errMsg := validation.ErrorMessage(validateErr)
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(errMsg))

			return
		}

		enrichPerson := enrichPerson(&person)

		id, err := creator.Create(r.Context(), enrichPerson)

		if errors.Is(err, repository.ErrPersonExists) {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, response.Error(fmt.Sprintf(
				"Person with name %s %s already exists", person.Name, person.Surname)))
			return
		}

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add person"))
			return
		}

		createResp := &models.Identifier{ID: id}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.OK[models.Identifier](createResp))
	}
}

func enrichPerson(person *models.Person) *models.EnrichmentPerson {
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
			log.Printf("failed to parse url for age: %v", err)
			return
		}

		resp, err := http.Get(fulUrlForAge)
		if err != nil {
			log.Printf("failed to fetch url for age: %v", err)
			return
		}
		defer resp.Body.Close()

		var data struct {
			Age *int64 `json:"age"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("failed to decode response body: %v", err)
			return
		}

		mu.Lock()
		enrichPerson.Age = data.Age
		mu.Unlock()
	}()

	go func() {
		defer wg.Done()

		fulUrlForGender, err := urlbuilder.BuildWithQueryParams(urlForGender, map[string]string{nameParam: person.Name})
		if err != nil {
			log.Printf("failed to parse url for gender: %v", err)
		}

		resp, err := http.Get(fulUrlForGender)
		if err != nil {
			log.Printf("failed to fetch url for gender: %v", err)
		}
		defer resp.Body.Close()

		var data struct {
			Gender      *models.Gender `json:"gender"`
			Probability float64        `json:"probability"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("failed to decode response body: %v", err)
		}

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
			log.Printf("failed to parse url for nationality: %v", err)
		}

		resp, err := http.Get(fullUrlNationality)
		if err != nil {
			log.Printf("failed to fetch url for nationality: %v", err)
		}
		defer resp.Body.Close()

		var data struct {
			Country []struct {
				CountryID   string  `json:"country_id"`
				Probability float64 `json:"probability"`
			} `json:"country"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Printf("failed to decode response body: %v", err)
		}

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
