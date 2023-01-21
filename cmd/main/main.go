package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-vendors-api/pkg/db"
)

var dbClient *db.DBClient

func main() {
	client, dbErr := db.InitDB()
	if dbErr != nil {
		panic(dbErr)
	}
	dbClient = &client

	router := gin.Default()

	// grouped routes for graphql playground and requests
	gql := router.Group("/graphql")
	{
		gql.GET("/", PlayGroundHandler())
		gql.POST("/query", GraphqlHandler())
	}

	router.GET("/vendors", GetAllVendors)
	router.GET("/vendors/:name", GetVendor)
	router.POST("/vendors", CreateVendor)
	router.PATCH("/vendors", UpdateVendor)
	router.DELETE("/vendors", DeleteVendor)
	router.Run() //runs on port 8080 by default
}
