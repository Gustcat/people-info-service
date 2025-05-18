package models

type Person struct {
	Name       string  `db:"name" json:"name" validate:"required,min=2,max=100"`
	Surname    string  `db:"surname" json:"surname" validate:"required,min=2,max=100"`
	Patronymic *string `db:"patronymic" json:"patronymic" validate:"min=2,max=100"`
}

type EnrichmentPerson struct {
	Person
	Age         *int64  `db:"age" json:"age" validate:"gte=0,lte=130"`
	Gender      *Gender `db:"gender" json:"gender" validate:"oneof=male female"`
	Nationality *string `db:"nationality" json:"nationality" validate:"min=2,max=100"`
}

type Identifier struct {
	ID int64 `db:"id" json:"id"`
}

type FullPerson struct {
	Identifier
	EnrichmentPerson
}

type PersonUpdate struct {
	Name        *string `db:"name" json:"name" validate:"min=2,max=100"`
	Surname     *string `db:"surname" json:"surname" validate:"min=2,max=100"`
	Patronymic  *string `db:"patronymic" json:"patronymic" validate:"min=2,max=100"`
	Age         *int64  `db:"age" json:"age" validate:"gte=0,lte=130"`
	Gender      *Gender `db:"gender" json:"gender" validate:"oneof=male female"`
	Nationality *string `db:"nationality" json:"nationality" validate:"min=2,max=100"`
}
