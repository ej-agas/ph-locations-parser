package psgc

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

func NewRow(row []string) *Row {

	status := ""
	if len(row) == 13 {
		status = row[12]
	}

	return &Row{
		PSGC:               row[0],
		Name:               row[1],
		CorrespondenceCode: row[2],
		GeographicLevel:    row[3],
		OldName:            row[4],
		CityClass:          row[5],
		IncomeClass:        row[6],
		UrbanRural:         row[7],
		Population2015:     row[8],
		Population2020:     row[10],
		Status:             status,
	}
}
