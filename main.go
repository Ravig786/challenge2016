package main

import (
	"fmt"

	"github.com/Ravig786/challenge2016/models"
	"github.com/Ravig786/challenge2016/router"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Initializing data...")

	models.InitRegionData()
	models.InitDistributorRegistry()

	err := models.LoadRegionDataFromCSV("cities.csv")
	if err != nil {
		fmt.Println("Error loading region data:", err)
	}
	// models.PrintGlobalRegionData()

	fmt.Println("Data initialized successfully!")

	r := gin.Default()
	router.SetupRoutes(r)
	r.Run(":5000")
}
