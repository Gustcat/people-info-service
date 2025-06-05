package models

type Person struct {
	Name       string  `db:"name" json:"name" binding:"required,min=2,max=100"`
	Surname    string  `db:"surname" json:"surname" binding:"required,min=2,max=100"`
	Patronymic *string `db:"patronymic" json:"patronymic" binding:"omitempty,min=2,max=100"`
}

type EnrichmentPerson struct {
	Person
	Age         *int64  `db:"age" json:"age" binding:"omitempty,gte=0,lte=130"`
	Gender      *Gender `db:"gender" json:"gender" binding:"omitempty,oneof=male female"`
	Nationality *string `db:"nationality" json:"nationality" binding:"omitempty,min=2,max=100"`
}

type Identifier struct {
	ID int64 `db:"id" json:"id"`
}

type FullPerson struct {
	Identifier
	EnrichmentPerson
}

type PersonUpdate struct {
	Name        *string `db:"name" json:"name" binding:"omitempty,min=2,max=100"`
	Surname     *string `db:"surname" json:"surname" binding:"omitempty,min=2,max=100"`
	Patronymic  *string `db:"patronymic" json:"patronymic" binding:"omitempty,min=2,max=100"`
	Age         *int64  `db:"age" json:"age" binding:"omitempty,gte=0,lte=130"`
	Gender      *Gender `db:"gender" json:"gender" binding:"omitempty,oneof=male female"`
	Nationality *string `db:"nationality" json:"nationality" binding:"omitempty,min=2,max=100"`
}
