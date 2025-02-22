package router

import (
	"github.com/Ravig786/challenge2016/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {

	// store distributor details
	r.POST("/distributor", handlers.CreateDistributorHandler)

	// give permission to distributor for location
	r.POST("/distributor/:name/:action", handlers.AddPermission)

	// check location for distributor
	r.GET("/distributor/:name/can-distribute", handlers.CheckDistribution)

	// get distributors
	r.GET("/distributors", handlers.GetAllDistributors)

	// get countries, states, cities
	r.GET("/countries", handlers.GetAllCountriesHandler)
	r.GET("/countries/:country_code/states", handlers.GetStatesByCountryHandler)
	r.GET("/countries/:country_code/states/:state_code/cities", handlers.GetCitiesByStateHandler)

}
