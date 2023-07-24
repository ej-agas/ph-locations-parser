package psgc

import (
	"context"
	"fmt"
	"github.com/ej-agas/ph-locations/models"
)

var (
	firstDistrict  = "NCR, City of Manila, First District (Not a Province)"
	secondDistrict = "NCR, Second District (Not a Province)"
	thirdDistrict  = "NCR, Third District (Not a Province)"
	fourthDistrict = "NCR, Fourth District (Not a Province)"
)

type Parser struct {
	State State
	Store Store
	Rows  []Row
	//CurrentRow chan CurrentRow
}

type CurrentRow struct {
	Type  string
	Name  string
	Count int
}

func (p Parser) Run(currentRow chan<- CurrentRow) error {
	var err error
	count := 1

	for _, row := range p.Rows {
		switch row.GeographicLevel {
		case "Reg":

			region := models.Region{
				Code:       row.PSGC,
				Name:       row.Name,
				Population: row.Population2020,
			}

			if err := p.Store.Region.Save(context.Background(), region); err != nil {
				return fmt.Errorf("error saving region: %s", err)
			}

			p.State.Region, err = p.Store.Region.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("failed to find region: %s", err)
			}

			currentRow <- CurrentRow{
				Type:  "Region",
				Name:  row.Name,
				Count: count,
			}
		case "Prov":

			province := models.Province{
				Code:        row.PSGC,
				Name:        row.Name,
				IncomeClass: row.IncomeClass,
				Population:  row.Population2020,
				RegionId:    &p.State.Region.Id,
			}

			if err := p.Store.Province.Save(context.Background(), province); err != nil {
				return fmt.Errorf("error saving province: %s", err)
			}

			p.State.Province, err = p.Store.Province.FindByCode(context.Background(), row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding province: %s", err)
			}

			currentRow <- CurrentRow{
				Type:  "Province",
				Name:  row.Name,
				Count: count,
			}
		case "Dist":

			district := models.District{
				Name:       row.Name,
				Population: row.Population2020,
				RegionId:   &p.State.Region.Id,
			}

			if row.Name == firstDistrict {
				district.Code = "1st"
			}

			if row.Name == secondDistrict {
				district.Code = "2nd"
			}

			if row.Name == thirdDistrict {
				district.Code = "3rd"
			}

			if row.Name == fourthDistrict {
				district.Code = "4th"
			}

			if err := p.Store.District.Save(context.Background(), district); err != nil {
				return fmt.Errorf("error saving district: %w", err)
			}

			p.State.District, err = p.Store.District.FindByCode(district.Code)
			if err != nil {
				return fmt.Errorf("error finding district: %w", err)
			}

			currentRow <- CurrentRow{
				Type:  "District",
				Name:  row.Name,
				Count: count,
			}
		case "City":
			city := models.City{
				Code:        row.PSGC,
				Name:        row.Name,
				CityClass:   row.CityClass,
				IncomeClass: row.IncomeClass,
				Population:  row.Population2020,
			}

			parentCode, err := NewPSGC(city.Code)

			if err != nil {
				fmt.Println(fmt.Errorf("invalid PSGC: %w", err))
			}

			if parentCode.Province() == p.State.Province.Code {
				city.ProvinceId = &p.State.Province.Id
			}

			if parentCode.Province() == p.State.District.Code {
				city.DistrictId = &p.State.District.Id
			}

			if err := p.Store.City.Save(context.Background(), city); err != nil {
				return fmt.Errorf("error saving city: %w", err)
			}

			p.State.City, err = p.Store.City.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding city: %w", err)

			}

			currentRow <- CurrentRow{
				Type:  "City",
				Name:  row.Name,
				Count: count,
			}
		case "Mun":
			municipality := models.Municipality{
				Code:        row.PSGC,
				Name:        row.Name,
				IncomeClass: row.IncomeClass,
				Population:  row.Population2020,
			}

			parentCode, err := NewPSGC(municipality.Code)

			if err != nil {
				return fmt.Errorf("invalid PSGC: %w", err)
			}

			if parentCode.Province() == p.State.Province.Code {
				municipality.ProvinceId = &p.State.Province.Id
			}

			if parentCode.Province() == p.State.District.Code {
				municipality.DistrictId = &p.State.District.Id
			}

			if err := p.Store.Municipality.Save(context.Background(), municipality); err != nil {
				return fmt.Errorf("error saving municipality: %w", err)
			}

			p.State.Municipality, err = p.Store.Municipality.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding municipality: %w", err)
			}

			currentRow <- CurrentRow{
				Type:  "Municipality",
				Name:  row.Name,
				Count: count,
			}
		case "SubMun":
			subMunicipality := models.SubMunicipality{
				Code:       row.PSGC,
				Name:       row.Name,
				Population: row.Population2020,
				CityId:     &p.State.City.Id,
			}

			if err := p.Store.SubMunicipality.Save(context.Background(), subMunicipality); err != nil {
				return fmt.Errorf("error saving sub municipality: %w", err)
			}

			p.State.SubMunicipality, err = p.Store.SubMunicipality.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding sub municipality %w", err)
			}

			currentRow <- CurrentRow{
				Type:  "Sub Municipality",
				Name:  row.Name,
				Count: count,
			}
		case "Bgy":
			barangay := models.Barangay{
				Code:       row.PSGC,
				Name:       row.Name,
				UrbanRural: row.UrbanRural,
				Population: row.Population2020,
			}

			parentCode, err := NewPSGC(barangay.Code)
			if err != nil {
				return fmt.Errorf("invalid PSGC: %w", err)
			}

			if parentCode.CityOrMunicipality() == p.State.City.Code {
				barangay.CityId = &p.State.City.Id
			}

			if parentCode.CityOrMunicipality() == p.State.Municipality.Code {
				barangay.MunicipalityId = &p.State.Municipality.Id
			}

			if parentCode.CityOrMunicipality() == p.State.SubMunicipality.Code {
				barangay.SubMunicipalityId = &p.State.SubMunicipality.Id
			}

			if parentCode.CityOrMunicipality() == p.State.SpecialGovernmentUnit.Code {
				barangay.SpecialGovernmentUnitId = &p.State.SpecialGovernmentUnit.Id
			}

			if err := p.Store.Barangay.Save(context.Background(), barangay); err != nil {
				return fmt.Errorf("error saving barangay %s: %w", barangay.Name, err)
			}

			currentRow <- CurrentRow{
				Type:  "Barangay",
				Name:  row.Name,
				Count: count,
			}
		case "SGU":
			sgu := models.SpecialGovernmentUnit{
				Code:       row.PSGC,
				Name:       row.Name,
				ProvinceId: &p.State.Province.Id,
			}

			if err := p.Store.SpecialGovernmentUnit.Save(context.Background(), sgu); err != nil {
				return fmt.Errorf("error saving special government unit %s: %w", sgu.Name, err)
			}

			p.State.SpecialGovernmentUnit, err = p.Store.SpecialGovernmentUnit.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding special government unit: %w", err)
			}

			currentRow <- CurrentRow{
				Type:  "Special Government Unit",
				Name:  row.Name,
				Count: count,
			}
		case "":
			province := models.Province{
				Code:     row.PSGC,
				Name:     row.Name,
				RegionId: &p.State.Region.Id,
			}

			if err := p.Store.Province.Save(context.Background(), province); err != nil {
				return fmt.Errorf("error saving province %s: %w", province.Name, err)
			}

			p.State.Province, err = p.Store.Province.FindByCode(context.Background(), row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding province: %s", err)
			}

			currentRow <- CurrentRow{
				Type:  "Interim Province",
				Name:  row.Name,
				Count: count,
			}
		}
		count++
	}

	return nil
}
