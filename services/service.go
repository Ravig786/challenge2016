package services

import (
	"fmt"
	"sort"

	"github.com/Ravig786/challenge2016/models"
	"github.com/Ravig786/challenge2016/utils"
)

func CreateDistributor(name, parent string) error {
	models.Registry.Lock()
	defer models.Registry.Unlock()

	if _, exists := models.Registry.Distributors[name]; exists {
		return fmt.Errorf("distributor %s already exists", name)
	}

	if parent == "" {
		models.Registry.Distributors[name] = &models.DistributorPermissions{
			Includes: make(map[string]struct{}),
			Excludes: make(map[string]struct{}),
			Parent:   "",
		}
		fmt.Printf("distributor %s created successfully\n", name)
		return nil
	}

	parentDistributor, parentExists := models.Registry.Distributors[parent]
	if !parentExists {
		return fmt.Errorf("parent distributor %s not found", parent)
	}

	subDistributor := &models.DistributorPermissions{
		Includes: make(map[string]struct{}),
		Excludes: make(map[string]struct{}),
		Parent:   parent,
	}

	for region := range parentDistributor.Includes {
		subDistributor.Includes[region] = struct{}{}
	}
	for region := range parentDistributor.Excludes {
		subDistributor.Excludes[region] = struct{}{}
	}

	models.Registry.Distributors[name] = subDistributor

	fmt.Printf("✅ Sub-distributor %s created under %s with inherited permissions\n", name, parent)
	return nil
}

func AddPermission(name, region string, isInclude bool) error {
	models.Registry.Lock()
	defer models.Registry.Unlock()

	distributor, exists := models.Registry.Distributors[name]
	if !exists {
		return fmt.Errorf("distributor %s not found", name)
	}

	normalizedRegion := utils.NormalizeRegion(region)

	fmt.Printf("Adding %s: %s for distributor: %s\n", map[bool]string{true: "INCLUDE", false: "EXCLUDE"}[isInclude], normalizedRegion, name)

	if !isRegionValid(normalizedRegion) {
		return fmt.Errorf("invalid region: %s. Region does not exist in memory", region)
	}

	if isInclude {
		distributor.Includes[normalizedRegion] = struct{}{}
	} else {
		distributor.Excludes[normalizedRegion] = struct{}{}
	}

	fmt.Printf("Updated Distributor Data: %+v\n", distributor)
	return nil
}

func CanDistribute(name, location string) (bool, error) {
	models.Registry.RLock()
	defer models.Registry.RUnlock()

	_, exists := models.Registry.Distributors[name]
	if !exists {
		return false, fmt.Errorf("distributor %s not found", name)
	}

	normalizedLocation := utils.NormalizeRegion(location)

	fmt.Printf("🔍 Checking if %s can distribute in: %s\n", name, normalizedLocation)

	if !isRegionValid(normalizedLocation) {
		fmt.Printf("⚠️ Invalid location check: %s does not exist in GlobalRegionData\n", location)
		return false, nil
	}

	locationParts := utils.SplitRegion(normalizedLocation)
	locationDepth := len(locationParts)

	for d := name; d != ""; d = models.Registry.Distributors[d].Parent {
		permissions, exists := models.Registry.Distributors[d]
		if !exists {
			return false, nil
		}

		if locationDepth == 3 {
			stateRegion := fmt.Sprintf("%s-%s", locationParts[1], locationParts[2])
			countryRegion := locationParts[2]

			if _, found := permissions.Excludes[countryRegion]; found {
				fmt.Printf("Access Denied! Country %s is excluded for %s\n", countryRegion, name)
				return false, nil
			}

			if _, found := permissions.Excludes[stateRegion]; found {
				fmt.Printf("Access Denied! State %s is excluded for %s\n", stateRegion, name)
				return false, nil
			}

			if _, found := permissions.Excludes[normalizedLocation]; found {
				fmt.Printf("Access Denied! City %s is explicitly excluded for %s\n", normalizedLocation, name)
				return false, nil
			}
		}

		if locationDepth == 2 {
			countryRegion := locationParts[1]

			if _, found := permissions.Excludes[countryRegion]; found {
				fmt.Printf("Access Denied! Country %s is excluded for %s\n", countryRegion, name)
				return false, nil
			}

			if _, found := permissions.Excludes[normalizedLocation]; found {
				fmt.Printf("Access Denied! State %s is explicitly excluded for %s\n", normalizedLocation, name)
				return false, nil
			}
		}

		if locationDepth == 1 {
			if _, found := permissions.Excludes[normalizedLocation]; found {
				fmt.Printf("Access Denied! Country %s is explicitly excluded for %s\n", normalizedLocation, name)
				return false, nil
			}
		}

		if locationDepth == 3 {
			stateRegion := fmt.Sprintf("%s-%s", locationParts[1], locationParts[2])
			countryRegion := locationParts[2]

			if _, found := permissions.Includes[countryRegion]; found {
				fmt.Printf("Access Allowed! Country-level access granted for %s\n", normalizedLocation)
				return true, nil
			}

			if _, found := permissions.Includes[stateRegion]; found {
				fmt.Printf("Access Allowed! State-level access granted for %s\n", normalizedLocation)
				return true, nil
			}

			if _, found := permissions.Includes[normalizedLocation]; found {
				fmt.Printf("Access Allowed! City %s is explicitly included for %s\n", normalizedLocation, name)
				return true, nil
			}
		}

		if locationDepth == 2 {
			countryRegion := locationParts[1]

			if _, found := permissions.Includes[countryRegion]; found {
				fmt.Printf("Access Allowed! Country-level access granted for %s\n", normalizedLocation)
				return true, nil
			}

			if _, found := permissions.Includes[normalizedLocation]; found {
				fmt.Printf("Access Allowed! State %s is explicitly included for %s\n", normalizedLocation, name)
				return true, nil
			}
		}

		if locationDepth == 1 {
			if _, found := permissions.Includes[normalizedLocation]; found {
				fmt.Printf("✅ Access Allowed! Country %s is explicitly included for %s\n", normalizedLocation, name)
				return true, nil
			}
		}
	}

	fmt.Printf("Access Denied! No INCLUDE rule found for %s\n", normalizedLocation)
	return false, nil
}

func GetAllDistributors() []string {
	var distributors []string
	for name := range models.Registry.Distributors {
		distributors = append(distributors, name)
	}
	return distributors
}

type countryInfo struct {
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
}

func GetAllCountries() []countryInfo {
	countries := []countryInfo{}

	for countryCode, country := range models.GlobalRegionData.Data {
		countries = append(countries, countryInfo{
			CountryCode: countryCode,
			CountryName: country.CountryName,
		})
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].CountryName < countries[j].CountryName
	})

	return countries
}

func GetStatesByCountry(countryCode string) ([]map[string]string, error) {

	fmt.Println(models.GlobalRegionData.Data)

	if country, exists := models.GlobalRegionData.Data[countryCode]; exists {
		var states []map[string]string
		for _, state := range country.States {
			states = append(states, map[string]string{
				"state_code": state.StateCode,
				"state_name": state.StateName,
			})
		}

		sort.Slice(states, func(i, j int) bool {
			return states[i]["state_name"] < states[j]["state_name"]
		})
		return states, nil
	}

	return nil, fmt.Errorf(" Country code %s not found", countryCode)
}

/*
func GetStatesByCountry(countryCode string) ([]models.State, error) {
	models.GlobalRegionData.RLock()
	defer models.GlobalRegionData.RUnlock()

	if country, exists := models.GlobalRegionData.Data[countryCode]; exists {
		return country.States, nil
	}

	return nil, fmt.Errorf("country code %s not found", countryCode)
}
*/

func GetCitiesByState(countryCode, stateCode string) ([]models.City, error) {

	if country, exists := models.GlobalRegionData.Data[countryCode]; exists {
		for _, state := range country.States {
			if state.StateCode == stateCode {
				sort.Slice(state.Cities, func(i, j int) bool {
					return state.Cities[i].CityName < state.Cities[j].CityName
				})
				return state.Cities, nil
			}
		}
		return nil, fmt.Errorf("state code %s not found in country %s", stateCode, countryCode)
	}

	return nil, fmt.Errorf("country code %s not found", countryCode)
}

func isRegionValid(region string) bool {

	parts := utils.SplitRegion(region)
	length := len(parts)

	switch length {
	case 1:
		_, exists := models.GlobalRegionData.Data[parts[0]]
		return exists
	case 2:
		if country, exists := models.GlobalRegionData.Data[parts[1]]; exists {
			for _, state := range country.States {
				if state.StateCode == parts[0] {
					return true
				}
			}
		}
	case 3:
		if country, exists := models.GlobalRegionData.Data[parts[2]]; exists {
			for _, state := range country.States {
				if state.StateCode == parts[1] {
					for _, city := range state.Cities {
						if city.CityCode == parts[0] {
							return true
						}
					}
				}
			}
		}
	}
	return false
}
