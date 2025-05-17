package models

import (
	"database/sql"
	"github.com/Gustcat/people-info-service/internal/models"
)

type Person struct {
	ID          int64          `db:"id" json:"id"`
	Name        string         `db:"name" json:"name"`
	Surname     string         `db:"surname" json:"surname"`
	Patronymic  sql.NullString `db:"patronymic" json:"patronymic,omitempty"`
	Age         int64          `db:"age" json:"age"`
	Gender      models.Gender  `db:"gender" json:"gender"`
	Nationality string         `db:"nationality" json:"nationality"`
}
