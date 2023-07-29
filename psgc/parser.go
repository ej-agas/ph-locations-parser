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
}

func (p Parser) Run(currentRow chan<- struct{}) error {
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

			currentRow <- struct{}{}
		case "Prov":

			province := models.Province{
				Code:        row.PSGC,
				Name:        row.Name,
				IncomeClass: row.IncomeClass,
				Population:  row.Population2020,
				RegionCode:  &p.State.Region.Code,
			}

			if err := p.Store.Province.Save(context.Background(), province); err != nil {
				return fmt.Errorf("error saving province: %s", err)
			}

			p.State.Province, err = p.Store.Province.FindByCode(context.Background(), row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding province: %s", err)
			}

			currentRow <- struct{}{}
		case "Dist":

			district := models.District{
				Name:       row.Name,
				Population: row.Population2020,
				RegionCode: &p.State.Region.Code,
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

			currentRow <- struct{}{}
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

			// Interim District codes of the 4 NCR districts
			NCRInterimDistrictCodes := []string{"1st", "2nd", "3rd", "4th"}
			for _, code := range NCRInterimDistrictCodes {
				NCRRegionCode := "1300000000"
				if p.State.District.Code == code && p.State.Region.Code == NCRRegionCode {
					city.DistrictCode = &p.State.District.Code
				}
			}

			if parentCode.Province() == p.State.Province.Code {
				city.ProvinceCode = &p.State.Province.Code
			}

			if parentCode.Province() == p.State.District.Code {
				city.DistrictCode = &p.State.District.Code
			}

			if err := p.Store.City.Save(context.Background(), city); err != nil {
				return fmt.Errorf("error saving city %s: %w", city.Name, err)
			}

			p.State.City, err = p.Store.City.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding city: %w", err)

			}

			currentRow <- struct{}{}
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
				municipality.ProvinceCode = &p.State.Province.Code
			}

			if parentCode.Province() == p.State.District.Code {
				municipality.DistrictCode = &p.State.District.Code
			}

			if err := p.Store.Municipality.Save(context.Background(), municipality); err != nil {
				return fmt.Errorf("error saving municipality: %w", err)
			}

			p.State.Municipality, err = p.Store.Municipality.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding municipality: %w", err)
			}

			currentRow <- struct{}{}
		case "SubMun":
			subMunicipality := models.SubMunicipality{
				Code:       row.PSGC,
				Name:       row.Name,
				Population: row.Population2020,
				CityCode:   &p.State.City.Code,
			}

			if err := p.Store.SubMunicipality.Save(context.Background(), subMunicipality); err != nil {
				return fmt.Errorf("error saving sub municipality: %w", err)
			}

			p.State.SubMunicipality, err = p.Store.SubMunicipality.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding sub municipality %w", err)
			}

			currentRow <- struct{}{}
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
				barangay.CityCode = &p.State.City.Code
			}

			if parentCode.CityOrMunicipality() == p.State.Municipality.Code {
				barangay.MunicipalityCode = &p.State.Municipality.Code
			}

			if parentCode.CityOrMunicipality() == p.State.SubMunicipality.Code {
				barangay.SubMunicipalityCode = &p.State.SubMunicipality.Code
			}

			if parentCode.CityOrMunicipality() == p.State.SpecialGovernmentUnit.Code {
				barangay.SpecialGovernmentUnitCode = &p.State.SpecialGovernmentUnit.Code
			}

			if err := p.Store.Barangay.Save(context.Background(), barangay); err != nil {
				return fmt.Errorf("error saving barangay %s: %w", barangay.Name, err)
			}

			currentRow <- struct{}{}
		case "SGU":
			sgu := models.SpecialGovernmentUnit{
				Code:         row.PSGC,
				Name:         row.Name,
				ProvinceCode: &p.State.Province.Code,
			}

			if err := p.Store.SpecialGovernmentUnit.Save(context.Background(), sgu); err != nil {
				return fmt.Errorf("error saving special government unit %s: %w", sgu.Name, err)
			}

			p.State.SpecialGovernmentUnit, err = p.Store.SpecialGovernmentUnit.FindByCode(row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding special government unit: %w", err)
			}

			currentRow <- struct{}{}
		case "":
			province := models.Province{
				Code:       row.PSGC,
				Name:       row.Name,
				RegionCode: &p.State.Region.Code,
			}

			if err := p.Store.Province.Save(context.Background(), province); err != nil {
				return fmt.Errorf("error saving province %s: %w", province.Name, err)
			}

			p.State.Province, err = p.Store.Province.FindByCode(context.Background(), row.PSGC)
			if err != nil {
				return fmt.Errorf("error finding province: %s", err)
			}

			currentRow <- struct{}{}
		}
		count++
	}

	return nil
}
