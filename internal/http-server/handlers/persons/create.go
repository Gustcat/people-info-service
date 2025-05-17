package persons

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/lib/urlbuilder"
	"github.com/Gustcat/people-info-service/internal/models"
	"github.com/go-chi/render"
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

type CreateResponse struct {
	ID int64 `json:"id"`
}

type Creator interface {
	Create(ctx context.Context, person *models.FullPerson) (int64, error)
}

func Create(ctx context.Context, creator Creator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.Create"

		var person models.FullPerson

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

		//TODO: обработка ошибок валидации

		enrichPerson(&person)

		id, err := creator.Create(r.Context(), &person)

		// TODO: обработать случай с попыткой создать уже существ. запись

		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, response.Error("failed to add person"))
			return
		}

		createResp := &CreateResponse{ID: id}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, response.OK[CreateResponse](createResp))
	}
}

func enrichPerson(person *models.FullPerson) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}

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
		person.Age = data.Age
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
		person.Gender = data.Gender
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
		person.Nationality = nationality
		mu.Unlock()
	}()

	wg.Wait()
}
