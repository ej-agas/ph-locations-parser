package cmd

import (
	"context"
	"fmt"
	"github.com/ej-agas/ph-locations/models"
	"github.com/ej-agas/ph-locations/postgresql"
	"github.com/ej-agas/psgc-publication-parser/psgc"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
	"strconv"
)

var (
	host     *string
	port     *string
	user     *string
	password *string
	db       *string
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse the PSGC .xlsx file",
	Long:  ``,
	Run:   process,
}

var (
	firstDistrict  = "NCR, City of Manila, First District (Not a Province)"
	secondDistrict = "NCR, Second District (Not a Province)"
	thirdDistrict  = "NCR, Third District (Not a Province)"
	fourthDistrict = "NCR, Fourth District (Not a Province)"
)

func foo(cmd *cobra.Command, args []string) {
	code, err := psgc.NewPSGC("1380614088s")

	if err != nil {
		fmt.Println(fmt.Errorf("invalid psgc code: %w", err))
		return
	}

	fmt.Printf("Region: %s\n", code.Region())
	fmt.Printf("Province: %s\n", code.Province())
	fmt.Printf("City or Municipality: %s\n", code.CityOrMunicipality())
	fmt.Printf("Barangay: %s\n", code.Barangay())
}

func process(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	dbPort, err := strconv.Atoi(*port)
	if err != nil {
		fmt.Println(fmt.Errorf("invalid port"))
		return
	}

	dbConfig := postgresql.Config{
		Host:         *host,
		Port:         dbPort,
		User:         *user,
		Password:     *password,
		DatabaseName: *db,
	}

	connection, err := postgresql.NewConnection(dbConfig)

	if err != nil {
		fmt.Println(fmt.Errorf("failed to connect to PostgreSQL: %s", err))
		return
	}

	regionStore := postgresql.NewRegionStore(connection)
	provinceStore := postgresql.NewProvinceStore(connection)
	districtStore := postgresql.NewDistrictStore(connection)
	cityStore := postgresql.NewCityStore(connection)
	municipalityStore := postgresql.NewMunicipalityStore(connection)
	subMunicipalityStore := postgresql.NewSubMunicipalityStore(connection)

	filePath := args[0]
	file, err := excelize.OpenFile(filePath)

	if err != nil {
		fmt.Println(fmt.Errorf("error opening file: %s", err))
		return
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(fmt.Errorf("error closing file: %s", err))
		}
	}()

	rows, err := file.GetRows("PSGC")
	if err != nil {
		fmt.Println(fmt.Errorf("error getting rows: %s", err))
		return
	}

	var currentRegion models.Region
	var currentProvince models.Province
	var currentDistrict models.District
	var currentCity models.City
	var currentMunicipality models.Municipality
	var currentSubMunicipality models.SubMunicipality

	rowCount := 1
	for _, item := range rows {
		// Skip the header row
		if rowCount == 1 {
			rowCount++
			continue
		}

		row := psgc.NewRow(item)

		switch row.GeographicLevel {
		case "Reg":

			region := models.Region{
				Code:       row.PSGC,
				Name:       row.Name,
				Population: row.IntPopulation2020(),
			}

			if err := regionStore.Save(context.Background(), region); err != nil {
				fmt.Println(fmt.Errorf("error saving region: %s", err))
				return
			}

			currentRegion, err = regionStore.FindByCode(row.PSGC)
			if err != nil {
				fmt.Println(fmt.Errorf("failed to find region: %s", err))
				return
			}

			fmt.Printf("saved region: %s\n", region.Name)
		case "Prov":

			province := models.Province{
				Code:        row.PSGC,
				Name:        row.Name,
				IncomeClass: row.IncomeClass,
				Population:  row.IntPopulation2020(),
				RegionId:    &currentRegion.Id,
			}

			if err := provinceStore.Save(context.Background(), province); err != nil {
				fmt.Println(fmt.Errorf("error saving province: %s", err))
				return
			}

			currentProvince, err = provinceStore.FindByCode(context.Background(), row.PSGC)
			if err != nil {
				fmt.Println(fmt.Errorf("error finding province: %s", err))
				return
			}

			fmt.Printf("saved province: %s\n", currentProvince.Name)
		case "Dist":

			district := models.District{
				Name:       row.Name,
				Population: row.IntPopulation2020(),
				RegionId:   &currentRegion.Id,
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

			if err := districtStore.Save(context.Background(), district); err != nil {
				fmt.Println(fmt.Errorf("error saving district: %w", err))
				return
			}

			currentDistrict, err = districtStore.FindByCode(district.Code)
			if err != nil {
				fmt.Println(fmt.Errorf("error finding district: %w", err))
				return
			}

			fmt.Printf("saved district: %s\n", currentDistrict.Name)
		case "City":
			city := models.City{
				Code:        row.PSGC,
				Name:        row.Name,
				CityClass:   row.CityClass,
				IncomeClass: row.IncomeClass,
				Population:  row.IntPopulation2020(),
			}

			parentCode, err := psgc.NewPSGC(city.Code)

			if err != nil {
				fmt.Println(fmt.Errorf("invalid PSGC: %w", err))
				return
			}

			if parentCode.Province() == currentProvince.Code {
				city.ProvinceId = &currentProvince.Id
			}

			if parentCode.Province() == currentDistrict.Code {
				city.DistrictId = &currentDistrict.Id
			}

			if err := cityStore.Save(context.Background(), city); err != nil {
				fmt.Println(fmt.Errorf("error saving city: %w", err))
			}

			currentCity, err = cityStore.FindByCode(row.PSGC)
			if err != nil {
				fmt.Println(fmt.Errorf("error finding city: %w", err))
				return
			}

			fmt.Printf("saved city: %s\n", currentCity.Name)
		case "Mun":
			municipality := models.Municipality{
				Code:        row.PSGC,
				Name:        row.Name,
				IncomeClass: row.IncomeClass,
				Population:  row.IntPopulation2020(),
			}

			parentCode, err := psgc.NewPSGC(municipality.Code)

			if err != nil {
				fmt.Println(fmt.Errorf("invalid PSGC: %w", err))
				return
			}

			if parentCode.Province() == currentProvince.Code {
				municipality.ProvinceId = &currentProvince.Id
			}

			if parentCode.Province() == currentDistrict.Code {
				municipality.DistrictId = &currentDistrict.Id
			}

			if err := municipalityStore.Save(context.Background(), municipality); err != nil {
				fmt.Printf("%#v\n", municipality)
				fmt.Println(fmt.Errorf("error saving municipality: %w", err))
			}

			currentMunicipality, err = municipalityStore.FindByCode(row.PSGC)
			if err != nil {
				fmt.Println(fmt.Errorf("error finding municipality: %w", err))
				return
			}

			fmt.Printf("saved municipality: %s\n", currentMunicipality.Name)
		case "SubMun":
			subMunicipality := models.SubMunicipality{
				Code:       row.PSGC,
				Name:       row.Name,
				Population: row.IntPopulation2020(),
				CityId:     &currentCity.Id,
			}

			if err := subMunicipalityStore.Save(context.Background(), subMunicipality); err != nil {
				fmt.Println(fmt.Errorf("error saving sub municipality: %w", err))
				return
			}

			currentSubMunicipality, err = subMunicipalityStore.FindByCode(row.PSGC)
			if err != nil {
				fmt.Println(fmt.Errorf("error finding sub municipality %w", err))
				return
			}

			fmt.Printf("saved sub municipality: %s\n", currentSubMunicipality.Name)
		case "Bgy":
			barangay := models.Barangay{
				Code:       row.PSGC,
				Name:       row.Name,
				UrbanRural: row.UrbanRural,
				Population: row.IntPopulation2020(),
			}

			parentCode, err := psgc.NewPSGC(barangay.Code)
			if err != nil {
				fmt.Println(fmt.Errorf("invalid PSGC: %w", err))
				return
			}

			if parentCode.CityOrMunicipality() == currentCity.Code {
				barangay.CityId = &currentCity.Id
			}

			if parentCode.CityOrMunicipality() == currentMunicipality.Code {
				barangay.MunicipalityId = &currentMunicipality.Id
			}

			if parentCode.CityOrMunicipality() == currentSubMunicipality.Code {
				barangay.SubMunicipalityId = &currentSubMunicipality.Id
			}

		case "":
		}
		rowCount++
	}
}

func init() {
	rootCmd.AddCommand(parseCmd)

	host = parseCmd.Flags().String("host", "127.0.0.1", "PostgreSQL Host")
	port = parseCmd.Flags().String("port", "5173", "PostgreSQL Port")
	user = parseCmd.Flags().String("user", "ph_locations_user", "PostgreSQL User")
	password = parseCmd.Flags().String("password", "", "PostgreSQL Password")
	db = parseCmd.Flags().String("db", "ph_locations_db", "PostgreSQL Database Name")
}
