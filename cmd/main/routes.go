package main

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-vendors-api/graph"
	"github.com/jhawk7/go-vendors-api/internal/handlers"
	"github.com/jhawk7/go-vendors-api/internal/pkg/db"
)

// Effectively sets up handler middleware for receiving and responding for graphql reqeusts
func PlayGroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL playground", "/graphql/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func GraphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// REST Handlers
func GetAllVendors(c *gin.Context) {
	vendors, err := dbClient.GetActiveVendors()
	if err != nil {
		handlers.ErrorHandler(c, err, http.StatusInternalServerError, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vendors,
	})
}

func GetVendor(c *gin.Context) {
	vendor, err, notFound := dbClient.GetVendorByName(c.Param("name"))
	if err != nil {
		var status int
		if notFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}
		handlers.ErrorHandler(c, err, status, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vendor,
	})
}

func CreateVendor(c *gin.Context) {
	var vendor db.Vendor
	if bindErr := c.Bind(&vendor); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	if createErr := dbClient.CreateVendor(&vendor); createErr != nil {
		handlers.ErrorHandler(c, createErr, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": vendor,
	})
}

func UpdateVendor(c *gin.Context) {
	var updateRequest db.UpdateRequest
	if bindErr := c.Bind(&updateRequest); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	vendor, updateErr, notFound := dbClient.UpdateVendor(updateRequest)
	if updateErr != nil {
		var status int
		if notFound {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}
		handlers.ErrorHandler(c, updateErr, status, false)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": vendor,
	})
}

func DeleteVendor(c *gin.Context) {
	var req db.DeleteRequest
	if bindErr := c.Bind(&req); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	dbClient.DeleteVendor(req.Name)
	c.Status(http.StatusNoContent)
}
