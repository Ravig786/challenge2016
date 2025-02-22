package handlers

import (
	"net/http"

	"github.com/Ravig786/challenge2016/services"
	"github.com/gin-gonic/gin"
)

func CreateDistributorHandler(c *gin.Context) {
	var payload struct {
		Name   string `json:"name"`
		Parent string `json:"parent,omitempty"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil || payload.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := services.CreateDistributor(payload.Name, payload.Parent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Distributor created", "name": payload.Name, "parent": payload.Parent})
}

func AddPermission(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")

	var payload struct {
		Region string `json:"region"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	isInclude := action == "include"
	err := services.AddPermission(name, payload.Region, isInclude)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission updated", "distributor": name, "action": action, "region": payload.Region})
}

func CheckDistribution(c *gin.Context) {
	name := c.Param("name")
	location := c.Query("location")

	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location query parameter is required"})
		return
	}

	allowed, err := services.CanDistribute(name, location)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"distributor": name, "location": location, "allowed": allowed})
}

func GetAllDistributors(c *gin.Context) {
	distributors := services.GetAllDistributors()
	c.JSON(http.StatusOK, gin.H{"distributors": distributors})
}

func GetAllCountriesHandler(c *gin.Context) {
	countries := services.GetAllCountries()
	c.JSON(http.StatusOK, gin.H{"countries": countries})
}

func GetStatesByCountryHandler(c *gin.Context) {
	countryCode := c.Param("country_code")

	states, err := services.GetStatesByCountry(countryCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"country_code": countryCode, "states": states})
}

func GetCitiesByStateHandler(c *gin.Context) {
	countryCode := c.Param("country_code")
	stateCode := c.Param("state_code")

	cities, err := services.GetCitiesByState(countryCode, stateCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"country_code": countryCode, "state_code": stateCode, "cities": cities})
}
