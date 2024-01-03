package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-vendors-api/graph"
	"github.com/jhawk7/go-vendors-api/internal/pkg/db"
	"github.com/jhawk7/go-vendors-api/internal/pkg/handlers"
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
func GetActiveVendors(c *gin.Context) {
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
	name := sanitizeName(c.Param("name"))
	vendor, notFound, err := dbClient.GetVendorByName(name)
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
	vendor := new(db.Vendor)
	if bindErr := c.Bind(vendor); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	vendor.Name = sanitizeName(vendor.Name)

	if createErr := dbClient.CreateVendor(vendor); createErr != nil {
		handlers.ErrorHandler(c, createErr, http.StatusBadRequest, false)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": vendor,
	})
}

func UpdateVendor(c *gin.Context) {
	updateRequest := new(db.UpdateRequest)
	if bindErr := c.Bind(updateRequest); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}
	updateRequest.Name = sanitizeName(updateRequest.Name)
	vendor, notFound, updateErr := dbClient.UpdateVendor(updateRequest)
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
	req := new(db.DeleteRequest)
	if bindErr := c.Bind(req); bindErr != nil {
		err := fmt.Errorf("failed to bind input params for request %v", bindErr)
		handlers.ErrorHandler(c, err, http.StatusBadRequest, false)
		return
	}

	req.Name = sanitizeName(req.Name)
	if err := dbClient.DeleteVendor(req.Name); err != nil {
		handlers.ErrorHandler(c, err, 0, false)
	}
	c.Status(http.StatusNoContent)
}

func sanitizeName(name string) string {
	safename := strings.TrimSpace(strings.ToLower(name))
	return safename
}
