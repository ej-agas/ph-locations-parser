package internal

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"regexp"
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
	Population2015     int
	Population2020     int
	Status             string
}

func NewRow(row []string) Row {
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
		Population2015:     0,
		Population2020:     0,
		Status:             "",
	}

	switch columnCount {
	case 0:
		return rowStruct
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
		rowStruct.Population2015 = strPopulationToInt(row[8])
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
		rowStruct.Population2015 = strPopulationToInt(row[8])
		rowStruct.Population2020 = strPopulationToInt(row[10])
	case 13:
		rowStruct.PSGC = row[0]
		rowStruct.Name = row[1]
		rowStruct.CorrespondenceCode = row[2]
		rowStruct.GeographicLevel = row[3]
		rowStruct.OldName = row[4]
		rowStruct.CityClass = row[5]
		rowStruct.IncomeClass = row[6]
		rowStruct.UrbanRural = row[7]
		rowStruct.Population2015 = strPopulationToInt(row[8])
		rowStruct.Population2020 = strPopulationToInt(row[10])
		rowStruct.Status = row[12]
	}

	return rowStruct
}

func strPopulationToInt(str string) int {
	re := regexp.MustCompile(`[0-9,]+`)
	match := re.FindString(str)

	match = regexp.MustCompile(`[^0-9]`).ReplaceAllString(match, "")
	num, err := strconv.Atoi(match)

	if err != nil {
		return 0
	}

	return num
}

func GetRowsFromFile(filePath string) ([]Row, error) {
	file, err := excelize.OpenFile(filePath)

	if err != nil {
		return make([]Row, 0, 0), fmt.Errorf("error opening file: %s", err)
	}

	defer file.Close()

	rawRows, err := file.GetRows("PSGC")
	if err != nil {
		return make([]Row, 0, 0), fmt.Errorf("error getting rows: %s", err)
	}

	rows := make([]Row, 0, len(rawRows)-1)

	isHeaderSkipped := false
	for _, rawRow := range rawRows {
		if isHeaderSkipped == false {
			isHeaderSkipped = true
			continue
		}

		rows = append(rows, NewRow(rawRow))
	}

	return rows, nil
}
