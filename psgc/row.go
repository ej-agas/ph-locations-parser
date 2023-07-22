package psgc

import (
	"strconv"
)

type Row struct {
	PSGC               string
	Name               string
	CorrespondenceCode string
	GeographicLevel    string
	OldName            string
	CityClass          string
	IncomeClass        string
	UrbanRural         string
	Population2015     string
	Population2020     string
	Status             string
}

func (r Row) IntPopulation2020() int {
	population, _ := strconv.Atoi(r.Population2020)

	return population
}

func NewRow(row []string) *Row {
	columnCount := len(row)

	rowStruct := Row{
		PSGC:               "",
		Name:               "",
		CorrespondenceCode: "",
		GeographicLevel:    "",
		OldName:            "",
		CityClass:          "",
		IncomeClass:        "",
		UrbanRural:         "",
		Population2015:     "",
		Population2020:     "",
		Status:             "",
	}

	switch columnCount {
	case 0:
		return &rowStruct
	case 1:
		rowStruct.PSGC = row[0]
	case 2:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
	case 3:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
	case 4:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
	case 5:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
	case 6:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]

	case 7:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
	case 8:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
		rowStruct.UrbanRural = row[7]
	case 9:
		fallthrough
	case 10:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
		rowStruct.UrbanRural = row[7]
		rowStruct.Population2015 = row[8]
	case 11:
		fallthrough
	case 12:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
		rowStruct.UrbanRural = row[7]
		rowStruct.Population2015 = row[8]
		rowStruct.Population2020 = row[10]
	case 13:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
		rowStruct.UrbanRural = row[7]
		rowStruct.Population2015 = row[8]
		rowStruct.Population2020 = row[10]
		rowStruct.Status = row[12]
	}

	return &rowStruct
}
