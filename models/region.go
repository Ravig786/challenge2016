package models

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
)

type City struct {
	CityCode string `json:"city_code"`
	CityName string `json:"city_name"`
}

type State struct {
	StateCode string `json:"state_code"`
	StateName string `json:"state_name"`
	Cities    []City `json:"city_list"`
}

type Country struct {
	CountryName string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	States      []State `json:"state_list"`
}

type RegionData struct {
	Data map[string]*Country
}

var GlobalRegionData *RegionData

func InitRegionData() {
	GlobalRegionData = &RegionData{
		Data: make(map[string]*Country),
	}
}

func LoadRegionDataFromCSV(filePath string) error {

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open region data file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, _ = reader.Read()

	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		cityCode, stateCode, countryCode := record[0], record[1], record[2]
		cityName, stateName, countryName := record[3], record[4], record[5]

		if _, exists := GlobalRegionData.Data[countryCode]; !exists {
			GlobalRegionData.Data[countryCode] = &Country{
				CountryCode: countryCode,
				CountryName: countryName,
				States:      []State{},
			}
		}

		country := GlobalRegionData.Data[countryCode]

		stateIndex := -1
		for i, state := range country.States {
			if state.StateCode == stateCode {
				stateIndex = i
				break
			}
		}
		if stateIndex == -1 {

			country.States = append(country.States, State{
				StateCode: stateCode,
				StateName: stateName,
				Cities:    []City{},
			})
			stateIndex = len(country.States) - 1
		}

		country.States[stateIndex].Cities = append(country.States[stateIndex].Cities, City{
			CityCode: cityCode,
			CityName: cityName,
		})
	}

	fmt.Println("Region data loaded successfully!")
	return nil
}

func PrintGlobalRegionData() {
	data, err := json.MarshalIndent(GlobalRegionData.Data, "", "  ")
	if err != nil {
		fmt.Println("error converting GlobalRegionData to JSON:", err)
		return
	}

	fmt.Println("globalRegionData (In-Memory Store):")
	fmt.Println(string(data))
}
