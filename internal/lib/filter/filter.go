package filter

type PersonFilter struct {
	Name        *string `schema:"name"`
	Surname     *string `schema:"surname"`
	Patronymic  *string `schema:"patronymic"`
	Gender      *string `schema:"gender"`
	AgeMin      *int64  `schema:"age_min"`
	AgeMax      *int64  `schema:"age_max"`
	Nationality *string `schema:"nationality"`
	Limit       *uint64 `schema:"limit"`
	Offset      *uint64 `schema:"offset"`
}
