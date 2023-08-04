package internal

import (
	"database/sql"
	"github.com/ej-agas/ph-locations/postgresql"
)

type Store struct {
	Region                *postgresql.RegionStore
	Province              *postgresql.ProvinceStore
	District              *postgresql.DistrictStore
	City                  *postgresql.CityStore
	Municipality          *postgresql.MunicipalityStore
	SubMunicipality       *postgresql.SubMunicipalityStore
	Barangay              *postgresql.BarangayStore
	SpecialGovernmentUnit *postgresql.SpecialGovernmentUnit
}

func NewStore(connection *sql.DB) *Store {
	return &Store{
		Region:                postgresql.NewRegionStore(connection),
		Province:              postgresql.NewProvinceStore(connection),
		District:              postgresql.NewDistrictStore(connection),
		City:                  postgresql.NewCityStore(connection),
		Municipality:          postgresql.NewMunicipalityStore(connection),
		SubMunicipality:       postgresql.NewSubMunicipalityStore(connection),
		Barangay:              postgresql.NewBarangayStore(connection),
		SpecialGovernmentUnit: postgresql.NewSpecialGovernmentUnit(connection),
	}
}
