package models

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

func IsValidGender(g Gender) bool {
	switch g {
	case GenderMale, GenderFemale:
		return true
	default:
		return false
	}
}
