package swagger

import (
	"github.com/Gustcat/people-info-service/internal/lib/response"
	"github.com/Gustcat/people-info-service/internal/models"
)

type FullPersonResponse struct {
	Status response.Status    `json:"status"  enums:"ok"`
	Data   *models.FullPerson `json:"data"`
}

type FullPersonsResponse struct {
	Status response.Status      `json:"status"  enums:"ok"`
	Data   []*models.FullPerson `json:"data"`
}

type IdResponse struct {
	Status response.Status   `json:"status"  enums:"ok"`
	Data   models.Identifier `json:"data"`
}

type EmptyResponse struct {
	Status response.Status `json:"status"  enums:"ok"`
}

type ErrorResponse struct {
	Status response.Status `json:"status" enums:"error" example:"error"`
	Error  string          `json:"error"`
}
