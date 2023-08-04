package internal

import "github.com/ej-agas/ph-locations/models"

type State struct {
	Region                models.Region
	Province              models.Province
	District              models.District
	City                  models.City
	Municipality          models.Municipality
	SubMunicipality       models.SubMunicipality
	Barangay              models.Barangay
	SpecialGovernmentUnit models.SpecialGovernmentUnit
}

func NewState() *State {
	return &State{}
}
